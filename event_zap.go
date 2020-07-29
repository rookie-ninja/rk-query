// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"reflect"
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

type Format int

const (
	JSON Format = 0
	RK   Format = 1
)

// Stringer above config file types.
func (fileType Format) String() string {
	names := [...]string{"JSON", "RK"}

	// Please do not forget to change the boundary while adding a new config file types
	if fileType < JSON || fileType > RK {
		return "UNKNOWN"
	}

	return names[fileType]
}

// It is not thread safe.
type EventZap struct {
	logger       *zap.Logger
	format       Format
	quietMode    bool
	appName      string
	hostname     string
	operation    string
	remoteAddr   string
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

func (event *EventZap) GetValue(key string) string {
	val, ok := event.pairs.Fields[key]
	str := cast.ToString(val)

	if ok && len(str) > 0 {
		return str
	} else {
		return ""
	}
}

func (event *EventZap) GetAppName() string {
	return event.appName
}

func (event *EventZap) GetHostname() string {
	return event.hostname
}

func (event *EventZap) GetLogger() *zap.Logger {
	return event.logger
}

func (event *EventZap) GetOperation() string {
	return event.operation
}

func (event *EventZap) SetOperation(operation string) {
	event.operation = operation
}

func (event *EventZap) GetEventStatus() eventStatus {
	return event.status
}

func (event *EventZap) SetStartTime(curr time.Time) {
	event.startTime = curr
	event.status = inProgress
}

func (event *EventZap) GetStartTime() time.Time {
	return event.startTime
}

func (event *EventZap) GetEndTime() time.Time {
	return event.endTime
}

func (event *EventZap) SetEndTime(curr time.Time) {
	if event.status != inProgress {
		return
	}

	event.endTime = curr

	if event.producesHistory() && event.eventHistory.builder.Len() > 0 {
		event.eventHistory.elapsedMS("end", toMillisecond(curr))
	}

	event.status = ended
}

func (event *EventZap) StartTimer(name string) {
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

func (event *EventZap) EndTimer(name string) {
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

func (event *EventZap) UpdateTimer(name string, elapsedMS int64) {
	event.UpdateTimerWithSample(name, elapsedMS, 1)
}

func (event *EventZap) UpdateTimerWithSample(name string, elapsedMS, sample int64) {
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

func (event *EventZap) GetTimeElapsedMS(name string) int64 {
	timer, contains := event.tracker[name]
	if !contains {
		return -1
	}

	return timer.GetElapsedMS()
}

func (event *EventZap) GetRemoteAddr() string {
	return event.remoteAddr
}

func (event *EventZap) SetRemoteAddr(addr string) {
	event.remoteAddr = addr
}

func (event *EventZap) GetCounter(key string) int64 {
	val, ok := event.counters.Fields[key]

	if ok {
		return cast.ToInt64(val)
	} else {
		return -1
	}
}

func (event *EventZap) SetCounter(key string, value int64) {
	event.counters.AddInt64(key, value)
}

func (event *EventZap) InCCounter(key string, delta int64) {
	val, ok := event.counters.Fields[key]

	if ok {
		event.counters.AddInt64(key, cast.ToInt64(val)+delta)
	} else {
		event.counters.AddInt64(key, delta)
	}
}

func (event *EventZap) AddPair(key, value string) {
	event.pairs.AddString(key, value)
}

func (event *EventZap) AddErr(err error) {
	name := reflect.TypeOf(err).Name()
	if len(name) < 1 {
		name = "std-err"
	}

	val, ok := event.errors.Fields[name]
	if !ok {
		event.errors.AddInt64(name, 1)
	} else {
		event.errors.AddInt64(name, cast.ToInt64(val)+1)
	}
}

func (event *EventZap) GetErrCount(err error) int64 {
	name := reflect.TypeOf(err).Name()
	if len(name) < 1 {
		name = "std-err"
	}

	val, ok := event.errors.Fields[name]
	if ok {
		return cast.ToInt64(val)
	}

	return 0
}

func (event *EventZap) AddFields(fields ...zap.Field) {
	event.fields = append(event.fields, fields...)
}

func (event *EventZap) GetFields() []zap.Field {
	return event.fields
}

func (event *EventZap) RecordHistoryEvent(name string) {
	if event.producesHistory() {
		event.eventHistory.elapsedMS(name, toMillisecond(time.Now()))
	}
}

func (event *EventZap) WriteLog() {
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

func (event *EventZap) toRkFormat() string {
	builder := &bytes.Buffer{}

	builder.WriteString(ScopeDelimiter + "\n")

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
	builder.WriteString(fmt.Sprintf("%s=%d\n", timeKey, event.GetEndTime().Sub(event.GetStartTime()).Milliseconds()))
	// hostname
	builder.WriteString(fmt.Sprintf("%s=%s\n", hostnameKey, event.GetHostname()))
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
	// operation
	builder.WriteString(fmt.Sprintf("%s=%s\n", operationKey, event.GetRemoteAddr()))
	// status
	builder.WriteString(fmt.Sprintf("%s=%s\n", eventStatusKey, event.GetEventStatus().String()))
	// history
	if event.producesHistory() && event.GetEventHistory().builder.Len() > 0 {
		builder.WriteString(historyKey + "=")
		event.GetEventHistory().appendTo(builder)
		builder.WriteString("\n")
	}

	builder.WriteString(EOE)
	return builder.String()
}

func (event *EventZap) toJsonFormat() []zap.Field {
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
	fields = append(fields, zap.Time(startTimeKey, event.GetStartTime()))
	// time
	fields = append(fields, zap.Int64(timeKey, event.GetEndTime().Sub(event.GetStartTime()).Milliseconds()))
	// hostname
	fields = append(fields, zap.String(hostnameKey, event.GetHostname()))
	// timing
	fields = append(fields, event.marshalTimerField())
	// counters
	fields = append(fields, zap.Any(counterKey, event.counters.Fields))
	// pairs
	fields = append(fields, zap.Any(pairKey, event.pairs.Fields))
	// err
	fields = append(fields, zap.Any(errKey, event.errors.Fields))
	// fields
	enc := zapcore.NewMapObjectEncoder()
	for _, v := range event.fields {
		v.AddTo(enc)
	}
	fields = append(fields, zap.Any(fieldKey, enc.Fields))

	// remote address
	fields = append(fields, zap.String(remoteAddrKey, event.GetRemoteAddr()))
	// app name
	fields = append(fields, zap.String(appNameKey, event.GetAppName()))

	// operation
	fields = append(fields, zap.String(operationKey, event.GetOperation()))

	// status
	fields = append(fields, zap.String(eventStatusKey, event.GetEventStatus().String()))

	// history
	if event.producesHistory() && event.GetEventHistory().builder.Len() > 0 {
		builder := &bytes.Buffer{}
		event.GetEventHistory().appendTo(builder)
		fields = append(fields, zap.String(historyKey, builder.String()))
	}

	return fields
}

func (event *EventZap) marshalEncoder(enc *zapcore.MapObjectEncoder) string {
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

func (event *EventZap) marshalTiming() string {
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

func (event *EventZap) marshalTimerField() zap.Field {
	enc := zapcore.NewMapObjectEncoder()
	for _, v := range event.tracker {
		v.ToZapFields(enc)
	}

	return zap.Any(timingKey, enc.Fields)
}

func (event *EventZap) GetEventHistory() *eventHistory {
	return event.eventHistory
}

func (event *EventZap) producesHistory() bool {
	return !event.quietMode
}

func (event *EventZap) inProgress() bool {
	if event.status != inProgress {
		return false
	}

	return true
}

func toMillisecond(curr time.Time) int64 {
	return curr.UnixNano() / 1e6
}
