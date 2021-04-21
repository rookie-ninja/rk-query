// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"time"
)

type eventStatus int

const (
	notStarted eventStatus = 0
	inProgress eventStatus = 1
	ended      eventStatus = 2
)

func (status eventStatus) String() string {
	names := [...]string{
		"NotStarted",
		"InProgress",
		"Ended",
	}

	if status < notStarted || status > ended {
		return "Unknown"
	}

	return names[status]
}

type format int

const (
	JSON format = 0
	RK   format = 1
)

// Stringer above config file types.
func (fileType format) String() string {
	names := [...]string{"JSON", "RK"}

	// Please do not forget to change the boundary while adding a new config file types
	if fileType < JSON || fileType > RK {
		return "UNKNOWN"
	}

	return names[fileType]
}

func ToFormat(f string) format {
	if f == "JSON" {
		return JSON
	} else if f == "RK" {
		return RK
	}

	return RK
}

// It is not thread safe.
type eventZap struct {
	logger       *zap.Logger
	format       format
	quietMode    bool
	appName      string
	appVersion   string
	locale       string
	entryName    string
	entryType    string
	hostname     string
	operation    string
	remoteAddr   string
	eventId      string
	resCode      string
	endTime      time.Time
	startTime    time.Time
	status       eventStatus
	counters     *zapcore.MapObjectEncoder
	pairs        *zapcore.MapObjectEncoder
	errors       *zapcore.MapObjectEncoder
	fields       []zap.Field
	eventHistory *eventHistory
	tracker      map[string]*timeTracker
}

func (event *eventZap) GetValue(key string) string {
	val, ok := event.pairs.Fields[key]
	str := cast.ToString(val)

	if ok && len(str) > 0 {
		return str
	} else {
		return ""
	}
}

func (event *eventZap) GetEntryName() string {
	return event.entryName
}

func (event *eventZap) GetEntryType() string {
	return event.entryType
}

func (event *eventZap) GetAppName() string {
	return event.appName
}

func (event *eventZap) GetAppVersion() string {
	return event.appVersion
}

func (event *eventZap) GetLocale() string {
	return event.locale
}

func (event *eventZap) GetEventId() string {
	return event.eventId
}

func (event *eventZap) SetEventId(id string) {
	event.eventId = id
}

func (event *eventZap) GetHostname() string {
	return event.hostname
}

func (event *eventZap) GetLogger() *zap.Logger {
	return event.logger
}

func (event *eventZap) GetOperation() string {
	return event.operation
}

func (event *eventZap) SetOperation(operation string) {
	event.operation = operation
}

func (event *eventZap) SetResCode(resCode string) {
	event.resCode = resCode
}

func (event *eventZap) GetEventStatus() eventStatus {
	return event.status
}

func (event *eventZap) SetStartTime(curr time.Time) {
	event.startTime = curr
	event.status = inProgress
}

func (event *eventZap) GetStartTime() time.Time {
	return event.startTime
}

func (event *eventZap) GetEndTime() time.Time {
	return event.endTime
}

func (event *eventZap) SetEndTime(curr time.Time) {
	if event.status != inProgress {
		return
	}

	event.endTime = curr

	if event.producesHistory() && event.eventHistory.builder.Len() > 0 {
		event.eventHistory.elapsedMS("end", toMillisecond(curr))
	}

	event.status = ended
}

func (event *eventZap) StartTimer(name string) {
	if !event.inProgress() || len(name) < 1 {
		return
	}

	_, contains := event.tracker[name]

	if !contains {
		tracker := newTimeTracker(name)
		if tracker == nil {
			return
		}

		event.tracker[name] = tracker
	}

	nowMS := toMillisecond(time.Now())
	tracker := event.tracker[name]
	tracker.Start(nowMS)

	if event.producesHistory() {
		event.eventHistory.elapsedMS("s-"+name, nowMS)
	}
}

func (event *eventZap) EndTimer(name string) {
	if !event.inProgress() || len(name) < 1 {
		return
	}

	tracker, contains := event.tracker[name]

	if !contains {
		return
	}

	nowMs := toMillisecond(time.Now())
	tracker.End(nowMs)

	if event.producesHistory() {
		event.eventHistory.elapsedMS("e-"+name, nowMs)
	}
}

func (event *eventZap) UpdateTimer(name string, elapsedMS int64) {
	event.UpdateTimerWithSample(name, elapsedMS, 1)
}

func (event *eventZap) UpdateTimerWithSample(name string, elapsedMS, sample int64) {
	if !event.inProgress() || len(name) < 1 {
		return
	}

	_, contains := event.tracker[name]

	if !contains {
		tracker := newTimeTracker(name)

		if tracker == nil {
			return
		}

		event.tracker[name] = tracker
	}

	tracker := event.tracker[name]
	tracker.ElapseWithSample(elapsedMS, sample)
}

func (event *eventZap) GetTimeElapsedMS(name string) int64 {
	timer, contains := event.tracker[name]
	if !contains {
		return -1
	}

	return timer.GetElapsedMS()
}

func (event *eventZap) GetRemoteAddr() string {
	return event.remoteAddr
}

func (event *eventZap) SetRemoteAddr(addr string) {
	event.remoteAddr = addr
}

func (event *eventZap) GetCounter(key string) int64 {
	val, ok := event.counters.Fields[key]

	if ok {
		return cast.ToInt64(val)
	} else {
		return -1
	}
}

func (event *eventZap) SetCounter(key string, value int64) {
	event.counters.AddInt64(key, value)
}

func (event *eventZap) InCCounter(key string, delta int64) {
	val, ok := event.counters.Fields[key]

	if ok {
		event.counters.AddInt64(key, cast.ToInt64(val)+delta)
	} else {
		event.counters.AddInt64(key, delta)
	}
}

func (event *eventZap) AddPair(key, value string) {
	event.pairs.AddString(key, value)
}

func (event *eventZap) AddErr(err error) {
	if err == nil {
		return
	}

	name := err.Error()

	if len(name) < 1 {
		name = "unknown"
	}

	val, ok := event.errors.Fields[name]
	if !ok {
		event.errors.AddInt64(name, 1)
	} else {
		event.errors.AddInt64(name, cast.ToInt64(val)+1)
	}
}

func (event *eventZap) GetErrCount(err error) int64 {
	name := err.Error()

	if len(name) < 1 {
		name = "unknown"
	}

	val, ok := event.errors.Fields[name]
	if ok {
		return cast.ToInt64(val)
	}

	return 0
}

func (event *eventZap) AddFields(fields ...zap.Field) {
	event.fields = append(event.fields, fields...)
}

func (event *eventZap) GetFields() []zap.Field {
	return event.fields
}

func (event *eventZap) RecordHistoryEvent(name string) {
	if event.producesHistory() {
		event.eventHistory.elapsedMS(name, toMillisecond(time.Now()))
	}
}

func (event *eventZap) WriteLog() {
	if event.format == JSON {
		event.GetLogger().With(event.toJsonFormat()...).Info("")
	} else {
		event.logger.Info(event.toRkFormat())
	}

	// finish any Time Aggregators that may not be done
	for _, v := range event.tracker {
		v.Finish()
	}
}

func (event *eventZap) toRkFormat() string {
	builder := &bytes.Buffer{}

	builder.WriteString(scopeDelimiter + "\n")

	// end_time
	if event.GetEndTime().IsZero() {
		event.SetEndTime(time.Now())
	}
	builder.WriteString(fmt.Sprintf("%s=%s\n", endTimeKey, event.GetEndTime().Format(time.RFC3339Nano)))
	// start_time
	if event.GetStartTime().IsZero() {
		event.SetStartTime(time.Now())
	}
	builder.WriteString(fmt.Sprintf("%s=%s\n", startTimeKey, event.GetStartTime().Format(time.RFC3339Nano)))
	// time
	builder.WriteString(fmt.Sprintf("%s=%d\n", elapsedKey, event.GetEndTime().Sub(event.GetStartTime()).Nanoseconds()))
	// hostname
	builder.WriteString(fmt.Sprintf("%s=%s\n", hostnameKey, event.GetHostname()))
	// eventId
	if len(event.GetEventId()) > 0 {
		builder.WriteString(fmt.Sprintf("%s=%s\n", eventIdKey, event.GetEventId()))
	}
	// timing
	builder.WriteString(fmt.Sprintf("%s=%s\n", timingKey, event.marshalTiming()))
	// counters
	builder.WriteString(fmt.Sprintf("%s=%s\n", counterKey, event.marshalEncoder(event.counters)))
	// pairs
	builder.WriteString(fmt.Sprintf("%s=%s\n", pairKey, event.marshalEncoder(event.pairs)))
	// errors
	builder.WriteString(fmt.Sprintf("%s=%s\n", errKey, event.marshalEncoder(event.errors)))
	// fields
	enc := zapcore.NewMapObjectEncoder()
	for i := range event.fields {
		event.fields[i].AddTo(enc)
	}
	bytes, _ := json.Marshal(enc.Fields)
	builder.WriteString(fmt.Sprintf("%s=%s\n", fieldKey, string(bytes)))
	// remote address
	builder.WriteString(fmt.Sprintf("%s=%s\n", remoteAddrKey, event.GetRemoteAddr()))

	// app name
	builder.WriteString(fmt.Sprintf("%s=%s\n", appNameKey, event.GetAppName()))
	// app version
	builder.WriteString(fmt.Sprintf("%s=%s\n", appVersionKey, event.GetAppVersion()))
	// entry name
	if len(event.GetEntryName()) > 0 {
		builder.WriteString(fmt.Sprintf("%s=%s\n", entryNameKey, event.GetEntryName()))
	}
	// entry type
	if len(event.GetEntryType()) > 0 {
		builder.WriteString(fmt.Sprintf("%s=%s\n", entryTypeKey, event.GetEntryType()))
	}
	// locale
	builder.WriteString(fmt.Sprintf("%s=%s\n", localeKey, event.GetLocale()))

	// operation
	builder.WriteString(fmt.Sprintf("%s=%s\n", operationKey, event.GetOperation()))
	// status
	builder.WriteString(fmt.Sprintf("%s=%s\n", eventStatusKey, event.GetEventStatus().String()))
	// resCode
	if len(event.resCode) > 0 {
		builder.WriteString(fmt.Sprintf("%s=%s\n", resCodeKey, event.resCode))
	}
	// history
	if event.producesHistory() && event.getEventHistory().builder.Len() > 0 {
		builder.WriteString(historyKey + "=")
		event.getEventHistory().appendTo(builder)
		builder.WriteString("\n")
	}

	// record timezone, os and arch
	zone, _ := time.Now().Zone()
	builder.WriteString(fmt.Sprintf("%s=%s\n", timezoneKey, zone))
	builder.WriteString(fmt.Sprintf("%s=%s\n", osKey, runtime.GOOS))
	builder.WriteString(fmt.Sprintf("%s=%s\n", archKey, runtime.GOARCH))

	builder.WriteString(eoe)
	return builder.String()
}

func (event *eventZap) toJsonFormat() []zap.Field {
	fields := make([]zapcore.Field, 0)

	// end_time
	if event.GetEndTime().IsZero() {
		event.SetEndTime(time.Now())
	}
	fields = append(fields, zap.Time(endTimeKey, event.GetEndTime()))
	// start_time
	if event.GetStartTime().IsZero() {
		event.SetStartTime(time.Now())
	}

	fields = append(fields,
		zap.Time(startTimeKey, event.GetStartTime()),
		zap.Int64(elapsedKey, event.GetEndTime().Sub(event.GetStartTime()).Nanoseconds()),
		zap.String(hostnameKey, event.GetHostname()),
		event.marshalTimerField(),
		zap.Any(counterKey, event.counters.Fields),
		zap.Any(pairKey, event.pairs.Fields),
		zap.Any(errKey, event.errors.Fields),
		zap.String(remoteAddrKey, event.GetRemoteAddr()),
		zap.String(appNameKey, event.GetAppName()),
		zap.String(appVersionKey, event.GetAppVersion()),
		zap.String(localeKey, event.GetLocale()),
		zap.String(operationKey, event.GetOperation()),
		zap.String(eventStatusKey, event.GetEventStatus().String()))

	// eventId
	if len(event.eventId) > 1 {
		fields = append(fields, zap.String(eventIdKey, event.GetEventId()))
	}
	// resCode
	if len(event.resCode) > 0 {
		fields = append(fields, zap.String(resCodeKey, event.resCode))
	}

	// fields
	enc := zapcore.NewMapObjectEncoder()
	for _, v := range event.fields {
		v.AddTo(enc)
	}
	fields = append(fields, zap.Any(fieldKey, enc.Fields))

	// entry name
	if len(event.GetEntryType()) > 0 {
		fields = append(fields, zap.String(entryNameKey, event.GetEntryName()))
	}

	// entry type
	if len(event.GetEntryType()) > 0 {
		fields = append(fields, zap.String(entryTypeKey, event.GetEntryType()))
	}

	// history
	if event.producesHistory() && event.getEventHistory().builder.Len() > 0 {
		builder := &bytes.Buffer{}
		event.getEventHistory().appendTo(builder)
		fields = append(fields, zap.String(historyKey, builder.String()))
	}

	// record timezone, os and arch
	zone, _ := time.Now().Zone()
	fields = append(fields,
		zap.String(timezoneKey, zone),
		zap.String(osKey, runtime.GOOS),
		zap.String(archKey, runtime.GOARCH))

	return fields
}

func (event *eventZap) marshalEncoder(enc *zapcore.MapObjectEncoder) string {
	builder := &bytes.Buffer{}

	bytes, err := json.Marshal(enc.Fields)
	if err != nil {
		builder.WriteString("{}")
	} else {
		_, err := builder.Write(bytes)
		if err != nil {
			builder.WriteString("{}")
		}
	}

	return builder.String()
}

func (event *eventZap) marshalTiming() string {
	builder := &bytes.Buffer{}
	enc := zapcore.NewMapObjectEncoder()

	for _, v := range event.tracker {
		v.ToZapFields(enc)
	}

	// all fields are int64
	bytes, err := json.Marshal(enc.Fields)
	if err != nil {
		builder.WriteString("{}")
	} else {
		builder.Write(bytes)
	}

	return builder.String()
}

func (event *eventZap) marshalTimerField() zap.Field {
	enc := zapcore.NewMapObjectEncoder()
	for _, v := range event.tracker {
		v.ToZapFields(enc)
	}

	return zap.Any(timingKey, enc.Fields)
}

func (event *eventZap) getEventHistory() *eventHistory {
	return event.eventHistory
}

func (event *eventZap) producesHistory() bool {
	return !event.quietMode
}

func (event *eventZap) inProgress() bool {
	if event.status != inProgress {
		return false
	}

	return true
}

func (event *eventZap) setLogger(logger *zap.Logger) {
	event.logger = logger
}

func (event *eventZap) setFormat(format format) {
	event.format = format
}

func (event *eventZap) setQuietMode(quietMode bool) {
	event.quietMode = quietMode
}

func (event *eventZap) setEntryName(entryName string) {
	event.entryName = entryName
}

func (event *eventZap) setEntryType(entryType string) {
	event.entryType = entryType
}

func (event *eventZap) setAppName(appName string) {
	event.appName = appName
}

func (event *eventZap) setAppVersion(appVersion string) {
	event.appVersion = appVersion
}

func (event *eventZap) setLocale(locale string) {
	event.locale = locale
}

func (event *eventZap) setHostname(hostname string) {
	event.hostname = hostname
}

func toMillisecond(curr time.Time) int64 {
	return curr.UnixNano() / 1e6
}
