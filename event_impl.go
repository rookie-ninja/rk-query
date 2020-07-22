// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// It is not thread safe.
type EventImpl struct {
	timeSource   TimeSource
	appName      string
	hostName     string
	operation    string
	remoteAddr   string
	zapLogger    *zap.Logger
	format       Format
	minimal      bool
	listeners    []eventEntryListener
	quietMode    bool
	defaultKvs   map[string]string
	mutex        *sync.Mutex
	eventHistory *eventHistory
	endTimeMS    int64
	startTimeMS  int64
	status       eventDataStatus
	counters     map[string]int64
	kvs          map[string]string
	errors       map[string]int64
	tracker      map[string]*timeTracker
}

func NewEventImpl(
	timeSource TimeSource,
	appName string,
	hostName string,
	defaultKvs map[string]string,
	zapLogger *zap.Logger,
	format Format,
	minimal bool,
	listeners []eventEntryListener,
	quietMode bool) *EventImpl {

	event := EventImpl{
		timeSource: timeSource,
		appName:    appName,
		hostName:   hostName,
		zapLogger:  zapLogger,
		listeners:  listeners,
		defaultKvs: defaultKvs,
		quietMode:  quietMode,
		format:     format,
		minimal:    minimal,
		tracker:    make(map[string]*timeTracker),
		kvs:        make(map[string]string),
		counters:   make(map[string]int64),
		errors:     make(map[string]int64),
		status:     notStarted,
		mutex:      &sync.Mutex{},
	}

	if !quietMode {
		event.eventHistory = newEventHistory()
	}

	return &event
}

func (event *EventImpl) GetAppName() string {
	return event.appName
}

func (event *EventImpl) GetHostName() string {
	return event.hostName
}

func (event *EventImpl) GetZapLogger() *zap.Logger {
	return event.zapLogger
}

func (event *EventImpl) GetOperation() string {
	return event.operation
}

func (event *EventImpl) SetOperation(operation string) {
	event.operation = operation
}

func (event *EventImpl) Reset() {
	event.status = notStarted
	event.startTimeMS = 0
	event.endTimeMS = 0
	event.counters = make(map[string]int64)
	event.kvs = make(map[string]string)
	event.tracker = make(map[string]*timeTracker)

	if event.eventHistory != nil {
		event.eventHistory.clear()
	}
}

func (event *EventImpl) GetEventStatus() eventDataStatus {
	return event.status
}

func (event *EventImpl) SetStartTimeMS(nowMS int64) {
	event.startTimeMS = nowMS
	event.status = inProgress
}

func (event *EventImpl) GetStartTimeMS() int64 {
	return event.startTimeMS
}

func (event *EventImpl) GetEndTimeMS() int64 {
	return event.endTimeMS
}

func (event *EventImpl) SetEndTimeMS(endTimeMS int64) {
	if event.status != inProgress {
		return
	}

	event.endTimeMS = endTimeMS

	if event.producesHistory() {
		event.eventHistory.elapsedMS("end", endTimeMS)
	}

	event.status = ended
}

func (event *EventImpl) StartTimer(name string) {
	if !event.inProgress() || len(name) < 1 {
		return
	}

	_, contains := event.tracker[name]

	if !contains {
		tracker := NewTimeTracker(name)
		if tracker == nil {
			return
		}

		event.tracker[name] = tracker
	}

	nowMS := event.timeSource.CurrentTimeMS()
	tracker := event.tracker[name]
	tracker.Start(nowMS)

	if event.producesHistory() {
		event.eventHistory.elapsedMS("s-"+name, nowMS)
	}
}

func (event *EventImpl) EndTimer(name string) {
	if !event.inProgress() {
		return
	}

	tracker, contains := event.tracker[name]

	if !contains {
		return
	}

	nowMS := event.timeSource.CurrentTimeMS()

	tracker.End(nowMS)

	if event.producesHistory() {
		event.eventHistory.elapsedMS("e-"+name, nowMS)
	}
}

func (event *EventImpl) UpdateTimer(name string, elapsedMS int64) {
	event.UpdateTimerWithSample(name, elapsedMS, 1)
}

func (event *EventImpl) UpdateTimerWithSample(name string, elapsedMS, sample int64) {
	if !event.inProgress() {
		return
	}

	_, contains := event.tracker[name]

	if !contains {
		tracker := NewTimeTracker(name)

		if tracker == nil {
			return
		}

		event.tracker[name] = tracker
	}

	tracker := event.tracker[name]
	tracker.ElapseWithSample(elapsedMS, sample)
}

func (event *EventImpl) GetTimeElapsedMS(timerName string) int64 {
	timer, contains := event.tracker[timerName]
	if !contains {
		return -1
	}

	return timer.GetElapsedMS()
}

func (event *EventImpl) GetRemoteAddr() string {
	return event.remoteAddr
}

func (event *EventImpl) SetRemoteAddr(addr string) {
	event.remoteAddr = addr
}

func (event *EventImpl) GetCounter(name string) int64 {
	value, contains := event.counters[name]
	if !contains {
		return -1
	}

	return value
}

func (event *EventImpl) SetCounter(name string, value int64) {
	event.counters[name] = value
}

func (event *EventImpl) InCCounter(name string, delta int64) {
	oldValue, contains := event.counters[name]

	if !contains {
		oldValue = delta
	} else {
		oldValue += delta
	}

	event.counters[name] = oldValue
}

func (event *EventImpl) AddKv(name, value string) {
	event.kvs[name] = value
}

func (event *EventImpl) AddErr(err error) {
	name := reflect.TypeOf(err).Name()
	if len(name) < 1 {
		name = "std-err"
	}

	existingValue, contains := event.errors[name]
	if !contains {
		event.errors[name] = 1
	} else {
		event.errors[name] = existingValue + 1
	}
}

func (event *EventImpl) AppendKv(name, value string) {
	existingValue, contains := event.kvs[name]
	if !contains {
		event.kvs[name] = value
	} else if len(existingValue) >= MaxHistoryLength {
		if !strings.HasSuffix(existingValue, CommaTruncated) {
			event.kvs[name] = existingValue + CommaTruncated
		}
	} else {
		event.kvs[name] = existingValue + "," + value
	}
}

func (event *EventImpl) GetValue(name string) string {
	value := event.kvs[name]
	return value
}

func (event *EventImpl) FinishCurrentEvent(name string) {
	if !event.inProgress() {
		return
	}

	nowMS := event.timeSource.CurrentTimeMS()

	event.UpdateTimer(name, nowMS-event.startTimeMS)

	if event.producesHistory() {
		event.eventHistory.elapsedMS(name, nowMS-event.startTimeMS)
	}
}

func (event *EventImpl) RecordHistoryEvent(name string) {
	if event.producesHistory() {
		event.eventHistory.elapsedMS(name, event.timeSource.CurrentTimeMS())
	}
}

func (event *EventImpl) WriteLog() {
	entry := eventEntryImpl{}

	if event.format == JSON {
		if event.minimal {
			entry.FormatAsJsonMin(event).Info("")
		} else {
			entry.FormatAsJson(event).Info("")
		}
	} else {
		if event.minimal {
			event.zapLogger.Info(event.toRkFormatMin())
		} else {
			event.zapLogger.Info(event.toRkFormat())
		}
	}

	// finish any Time Aggregators that may not be done
	for _, v := range event.tracker {
		v.Finish(event.timeSource)
	}

	for _, value := range event.listeners {
		value.notify(&entry)
	}
}

func (event *EventImpl) toRkFormat() string {
	entry := eventEntryImpl{}
	return entry.FormatAsRk(event)
}

func (event *EventImpl) toRkFormatMin() string {
	entry := eventEntryImpl{}
	return entry.FormatAsRkMin(event)
}

func (event *EventImpl) ToZapFields() []zap.Field {
	fields := make([]zapcore.Field, 0)

	entry := eventEntryImpl{}
	// EndTime
	fields = append(fields, zap.Time("end_time", time.Unix(0, event.GetEndTimeMS()*1000000)))
	// StartTime
	fields = append(fields, zap.Time("start_time", time.Unix(0, event.GetStartTimeMS()*1000000)))
	// Remote Address
	if len(event.GetRemoteAddr()) > 0 {
		fields = append(fields, zap.String("remote_address", event.GetRemoteAddr()))
	}
	// App
	if len(event.GetAppName()) > 0 {
		fields = append(fields, zap.String("app", event.GetAppName()))
	}
	// Hostname
	if len(event.GetHostName()) > 0 {
		fields = append(fields, zap.String("host_name", event.GetHostName()))
	}
	// Status
	if len(event.GetEventStatus().String()) > 0 {
		fields = append(fields, zap.String("event_status", event.GetEventStatus().String()))
	}
	// History
	if event.producesHistory() {
		builder := &bytes.Buffer{}
		event.GetEventHistory().appendTo(builder)

		fields = append(fields, zap.String("history", builder.String()))
	}
	// Timers
	fields = append(fields, entry.getTimerAsZapFields(event)...)
	// Counters
	fields = append(fields, entry.getCounterAsZapFields(event)...)
	// Name values
	fields = append(fields, entry.getKvAsZapFields(event)...)
	// err
	fields = append(fields, entry.getErrAsZapFields(event)...)
	return fields
}

func (event *EventImpl) ToZapFieldsMin() []zap.Field {
	fields := make([]zapcore.Field, 0)

	entry := eventEntryImpl{}

	// Timers
	fields = append(fields, entry.getTimerAsZapFields(event)...)
	// Counters
	fields = append(fields, entry.getCounterAsZapFields(event)...)
	// Kvs
	fields = append(fields, entry.getKvAsZapFields(event)...)
	// err
	fields = append(fields, entry.getErrAsZapFields(event)...)

	return fields
}

func (event *EventImpl) GetEventHistory() *eventHistory {
	return event.eventHistory
}

// Custom function
func (event *EventImpl) hasErrs() bool {
	return len(event.errors) > 0
}

func (event *EventImpl) hasCounters() bool {
	return len(event.counters) > 0
}

func (event *EventImpl) hasKvs() bool {
	return len(event.kvs) > 0
}

func (event *EventImpl) appendCounters(builder *bytes.Buffer) {
	var isFirst = true

	for k, v := range event.counters {
		if isFirst {
			isFirst = false
		} else {
			builder.WriteByte(',')
		}

		builder.WriteString(k + "=" + strconv.FormatInt(v, 10))
	}
}

func (event *EventImpl) appendErrs(builder *bytes.Buffer) {
	var isFirst = true

	for k, v := range event.errors {
		if isFirst {
			isFirst = false
		} else {
			builder.WriteByte(',')
		}

		builder.WriteString(k + "=" + strconv.FormatInt(v, 10))
	}
}

func (event *EventImpl) appendTimers(builder *bytes.Buffer) {
	var isFirst = true

	for _, v := range event.tracker {
		if isFirst {
			isFirst = false
		} else {
			builder.WriteByte(',')
		}
		str := v.StringWithTimeSource(event.timeSource)
		builder.WriteString(str)
	}
}

func (event *EventImpl) addDefaultKvs() {
	event.mutex.Lock()
	defer event.mutex.Unlock()

	for k, v := range event.defaultKvs {
		_, contains := event.kvs[k]
		if contains {
			continue
		}

		event.AddKv(k, v)
	}
}

func (event *EventImpl) appendKvs(builder *bytes.Buffer) {
	event.addDefaultKvs()

	var keys []string

	for k := range event.kvs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	if len(keys) > 0 {
		builder.WriteString(keys[0] + "=" + event.kvs[keys[0]])

	}

	for i := 1; i < len(keys); i++ {
		builder.WriteString("," + keys[i] + "=" + event.kvs[keys[i]])
	}
}

func (event *EventImpl) producesHistory() bool {
	return !event.quietMode
}

func (event *EventImpl) inProgress() bool {
	if event.status != inProgress {
		return false
	}

	return true
}
