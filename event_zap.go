// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkquery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"text/tabwriter"
	"time"
)

type eventStatus int

const (
	// NotStarted will be assigned to Event before SetStartTime() was called.
	NotStarted eventStatus = 0
	// InProgress will be assigned to Event after SetStartTime() was called.
	InProgress eventStatus = 1
	// Ended will be assigned to Event after SetEndTime() was called.
	Ended eventStatus = 2
)

// String will return string value of event status.
func (status eventStatus) String() string {
	names := [...]string{
		"NotStarted",
		"InProgress",
		"Ended",
	}

	if status < NotStarted || status > Ended {
		return "Unknown"
	}

	return names[status]
}

// Encoding supported format of Console and JSON currently.
type Encoding int

const (
	// CONSOLE is human readable format.
	CONSOLE Encoding = 0
	// JSON format.
	JSON Encoding = 1
	// FLATTEN format.
	FLATTEN Encoding = 2
)

// String will return string value of Encoding types.
func (ec Encoding) String() string {
	names := [...]string{"console", "json", "flatten"}

	// Please do not forget to change the boundary while adding a new config file types
	if ec > FLATTEN || ec < CONSOLE {
		return "UNKNOWN"
	}

	return names[ec]
}

// ToEncoding returns Encoding type from string value.
func ToEncoding(f string) Encoding {
	f = strings.ToLower(f)
	switch f {
	case "json":
		return JSON
	case "console":
		return CONSOLE
	case "flatten":
		return FLATTEN
	default:
		return CONSOLE
	}

	return CONSOLE
}

// It is not thread safe.
type eventZap struct {
	logger     *zap.Logger
	encoding   Encoding
	quietMode  bool
	appName    string                    // Application
	appVersion string                    // Application
	entryName  string                    // Application
	entryType  string                    // Application
	eventId    string                    // Ids
	traceId    string                    // Ids
	requestId  string                    // Ids
	endTime    time.Time                 // Time
	startTime  time.Time                 // Time
	timeZone   string                    // Time
	payloads   []zap.Field               // Payloads
	errors     *zapcore.MapObjectEncoder // Error
	operation  string                    // Event
	remoteAddr string                    // Event
	resCode    string                    // Event
	status     eventStatus               // Event
	pairs      *zapcore.MapObjectEncoder // Event
	counters   *zapcore.MapObjectEncoder // Event
	tracker    map[string]*timeTracker   // Event
}

// ************* Time *************

// SetStartTime sets start timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventZap) SetStartTime(curr time.Time) {
	event.startTime = curr
	event.status = InProgress
}

// GetStartTime Get start time of current event data.
func (event *eventZap) GetStartTime() time.Time {
	return event.startTime
}

// SetEndTime sets end timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventZap) SetEndTime(curr time.Time) {
	if event.status != InProgress {
		return
	}

	event.endTime = curr
	event.status = Ended
}

// GetEndTime returns end time of current event data.
func (event *eventZap) GetEndTime() time.Time {
	return event.endTime
}

// ************* Payload *************

// AddPayloads function add payload as zap.Field.
// Payload could be anything with RPC requests or user event such as http request param.
func (event *eventZap) AddPayloads(fields ...zap.Field) {
	event.payloads = append(event.payloads, fields...)
}

// ListPayloads will lists payloads.
func (event *eventZap) ListPayloads() []zap.Field {
	return event.payloads
}

// ************* Identity *************

// GetEventId returns event id of current event.
func (event *eventZap) GetEventId() string {
	return event.eventId
}

// SetEventId sets event id of current event.
// A new event id would be created while event data was created from EventFactory.
// User could override event id with this function.
func (event *eventZap) SetEventId(id string) {
	event.eventId = id
}

// GetTraceId returns trace id of current event.
func (event *eventZap) GetTraceId() string {
	return event.traceId
}

// SetTraceId set trace id of current event.
func (event *eventZap) SetTraceId(id string) {
	event.traceId = id
}

// GetRequestId returns request id of current event.
func (event *eventZap) GetRequestId() string {
	return event.requestId
}

// SetRequestId set request id of current event.
func (event *eventZap) SetRequestId(id string) {
	event.requestId = id
}

// ************* Error *************

// AddErr function adds an error into event which could be printed with error.Error() function.
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

// GetErrCount returns error count.
// We will use value of error.Error() as the key.
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

// ************* Event *************

// GetOperation returns operation of current event.
func (event *eventZap) GetOperation() string {
	return event.operation
}

// SetOperation sets operation of current event.
func (event *eventZap) SetOperation(operation string) {
	event.operation = operation
}

// GetRemoteAddr returns remote address of current event.
func (event *eventZap) GetRemoteAddr() string {
	return event.remoteAddr
}

// SetRemoteAddr sets remote address of current event, mainly used in RPC calls.
// Default value of <localhost> would be assigned while creating event via EventFactory.
func (event *eventZap) SetRemoteAddr(addr string) {
	event.remoteAddr = addr
}

// GetResCode returns response code of current event.
// Mainly used in RPC calls.
func (event *eventZap) GetResCode() string {
	return event.resCode
}

// SetResCode sets response code of current event.
func (event *eventZap) SetResCode(resCode string) {
	event.resCode = resCode
}

// GetEventStatus returns event status of current event.
// Available event status as bellow:
// 1: NotStarted
// 2: InProgress
// 3: Ended
func (event *eventZap) GetEventStatus() eventStatus {
	return event.status
}

// StartTimer starts timer of current sub event.
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

	nowMs := toMillisecond(time.Now())
	tracker := event.tracker[name]
	tracker.Start(nowMs)
}

// EndTimer ends timer of current sub event.
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
}

// UpdateTimerMs updates timer of current sub event with time elapsed in milli seconds.
func (event *eventZap) UpdateTimerMs(name string, elapsedMs int64) {
	event.UpdateTimerMsWithSample(name, elapsedMs, 1)
}

// UpdateTimerMsWithSample updates timer of current sub event with time elapsed in milli seconds.
func (event *eventZap) UpdateTimerMsWithSample(name string, elapsedMs, sample int64) {
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
	tracker.ElapseWithSample(elapsedMs, sample)
}

// GetTimeElapsedMs returns timer elapsed in milli seconds.
func (event *eventZap) GetTimeElapsedMs(name string) int64 {
	timer, contains := event.tracker[name]
	if !contains {
		return -1
	}

	return timer.GetElapsedMs()
}

// GetValueFromPair returns value with key in pairs.
func (event *eventZap) GetValueFromPair(key string) string {
	val, ok := event.pairs.Fields[key]
	str := cast.ToString(val)

	if ok && len(str) > 0 {
		return str
	}

	return ""
}

// AddPair adds value with key in pairs.
func (event *eventZap) AddPair(key, value string) {
	event.pairs.AddString(key, value)
}

// GetCounter returns counter of current event.
func (event *eventZap) GetCounter(key string) int64 {
	val, ok := event.counters.Fields[key]

	if ok {
		return cast.ToInt64(val)
	}

	return -1
}

// SetCounter sets counter of current event.
func (event *eventZap) SetCounter(key string, value int64) {
	event.counters.AddInt64(key, value)
}

// IncCounter increases counter of current event.
func (event *eventZap) IncCounter(key string, delta int64) {
	val, ok := event.counters.Fields[key]

	if ok {
		event.counters.AddInt64(key, cast.ToInt64(val)+delta)
	} else {
		event.counters.AddInt64(key, delta)
	}
}

// Finish sets event status and flush to logger.
func (event *eventZap) Finish() {
	if event.quietMode {
		return
	}

	switch event.encoding {
	case JSON:
		event.logger.With(event.toJsonFormat()...).Info("")
	case CONSOLE:
		event.logger.Info(event.toConsoleFormat())
	case FLATTEN:
		event.logger.Info(event.toFlattenFormat())
	default:
		event.logger.Info(event.toConsoleFormat())
	}

	// finish any Time Aggregators that may not be done
	for _, v := range event.tracker {
		v.Finish()
	}
}

// Sync flushes logs in buffer, mainly used for external syncer
func (event *eventZap) Sync() {
	event.logger.Sync()
}

// ************* Internal *************

// Marshal to FLATTEN format.
func (event *eventZap) toFlattenFormat() string {
	builder := &bytes.Buffer{}
	writer := tabwriter.NewWriter(builder, 2, 0, 4, ' ', tabwriter.TabIndent|tabwriter.StripEscape)

	// timestamp
	fmt.Fprint(writer, fmt.Sprintf("%s", event.GetEndTime().Format("2006-01-02T15:04:05.000Z0700")))

	// res code
	fmt.Fprint(writer, fmt.Sprintf("\t[%s]", getDefaultIfEmptyString(event.resCode, "[X]")))

	// elapsed
	fmt.Fprint(writer, fmt.Sprintf("\t%dms", event.GetEndTime().Sub(event.GetStartTime()).Milliseconds()))

	// API method
	// distinguish restful API and gRPC
	var grpcMethod, grpcServer, grpcType, apiPath, apiMethod, apiProtocol *zap.Field
	for i := range event.payloads {
		field := event.payloads[i]
		switch field.Key {
		case "grpcMethod":
			grpcMethod = &field
		case "grpcServer":
			grpcServer = &field
		case "grpcType":
			grpcType = &field
		case "apiPath":
			apiPath = &field
		case "apiMethod":
			apiMethod = &field
		case "apiProtocol":
			apiProtocol = &field
		}
	}

	method, operation, protocol := "", "", ""
	if grpcMethod != nil && grpcServer != nil {
		method = grpcServer.String
		operation = grpcMethod.String
		protocol = grpcType.String
	} else if apiPath != nil && apiMethod != nil {
		method = apiMethod.String
		operation = apiPath.String
		protocol = apiProtocol.String
	} else {
		operation = event.operation
		method = event.entryName
		protocol = event.entryType
	}

	// operation
	fmt.Fprint(writer, fmt.Sprintf("\t%s", getDefaultIfEmptyString(operation, "[X]")))

	// method
	fmt.Fprint(writer, fmt.Sprintf("\t%s", getDefaultIfEmptyString(method, "[X]")))

	// protocol
	fmt.Fprint(writer, fmt.Sprintf("\t%s", getDefaultIfEmptyString(protocol, "[X]")))

	// remote addr
	fmt.Fprint(writer, fmt.Sprintf("\t%s", getDefaultIfEmptyString(event.remoteAddr, "[X]")))

	// ids
	ids := make([]string, 0)
	if len(event.eventId) > 0 {
		ids = append(ids, event.eventId)
	}
	if len(event.traceId) > 0 {
		ids = append(ids, event.traceId)
	}
	fmt.Fprint(writer, fmt.Sprintf("\t[%s]", getDefaultIfEmptyString(strings.Join(ids, ","), "[X]")))

	writer.Flush()
	return builder.String()
}

// Marshal to CONSOLE format.
func (event *eventZap) toConsoleFormat() string {
	builder := &bytes.Buffer{}

	builder.WriteString(scopeDelimiter + "\n")

	// We would expect bellow format of event data as RK format.
	// ------------------------------------------------------------------------
	// endTime=2021-06-13T00:24:20.256315+08:00
	// startTime=2021-06-13T00:24:19.251056+08:00
	// elapsedNano=1005258286
	// timezone=CST
	// ids={"eventId":"6a2f84a8-a09a-42dc-bc9e-cabc7977345d"}
	// app={"appName":"appName","appVersion":"v0.0.1","entryName":"entry-example","entryType":"example"}
	// env={"arch":"amd64","hostname":"lark.local","localIP":"localhost","realm":"rk","region":"ap-guangzhou","az":"ap-guangzhou-1","domain":"beta","os":"darwin"}
	// payloads={"f1":"f2","t2":"2021-06-13T00:24:20.256276+08:00"}
	// error={"my error":1}
	// counters={"count":1}
	// pairs={"key":"value"}
	// timing={"t1.count":1,"t1.elapsed_ms":1005}
	// remoteAddr=localhost
	// operation=op
	// resCode=200
	// eventStatus=Ended
	// EOE

	// ************* Time *************
	// endTime
	if event.GetEndTime().IsZero() {
		event.SetEndTime(time.Now())
	}
	builder.WriteString(fmt.Sprintf("%s=%s\n", endTimeKey, event.GetEndTime().Format(time.RFC3339Nano)))
	// startTime
	if event.GetStartTime().IsZero() {
		event.SetStartTime(time.Now())
	}
	builder.WriteString(fmt.Sprintf("%s=%s\n", startTimeKey, event.GetStartTime().Format(time.RFC3339Nano)))
	// elapsedNano
	builder.WriteString(fmt.Sprintf("%s=%d\n", elapsedKey, event.GetEndTime().Sub(event.GetStartTime()).Nanoseconds()))
	// timeZone
	builder.WriteString(fmt.Sprintf("%s=%s\n", timezoneKey, event.timeZone))

	// ************* Ids *************
	builder.WriteString(fmt.Sprintf("%s=%s\n", idsKey, event.marshalIds()))

	// ************* App *************
	builder.WriteString(fmt.Sprintf("%s=%s\n", appKey, event.marshalApp()))

	// ************* Env *************
	builder.WriteString(fmt.Sprintf("%s=%s\n", envKey, event.marshalEnv()))

	// ************* Payloads *************
	builder.WriteString(fmt.Sprintf("%s=%s\n", payloadsKey, event.marshalPayloads()))

	// ************* Error *************
	if len(event.errors.Fields) > 0 {
		builder.WriteString(fmt.Sprintf("%s=%s\n", errKey, event.marshalEncoder(event.errors)))
	}

	// ************* Counter *************
	builder.WriteString(fmt.Sprintf("%s=%s\n", countersKey, event.marshalEncoder(event.counters)))

	// ************* Pairs *************
	builder.WriteString(fmt.Sprintf("%s=%s\n", pairsKey, event.marshalEncoder(event.pairs)))

	// ************* Timing *************
	builder.WriteString(fmt.Sprintf("%s=%s\n", timingKey, event.marshalTiming()))

	// ************* Event *************
	// remote address
	builder.WriteString(fmt.Sprintf("%s=%s\n", remoteAddrKey, event.GetRemoteAddr()))
	// operation
	builder.WriteString(fmt.Sprintf("%s=%s\n", operationKey, event.GetOperation()))
	// resCode
	if len(event.resCode) > 0 {
		builder.WriteString(fmt.Sprintf("%s=%s\n", resCodeKey, event.resCode))
	}
	// status
	builder.WriteString(fmt.Sprintf("%s=%s\n", eventStatusKey, event.GetEventStatus().String()))

	builder.WriteString(eoe)
	return builder.String()
}

// Marshal to JSON format.
func (event *eventZap) toJsonFormat() []zap.Field {
	fields := make([]zapcore.Field, 0)

	// We would expect bellow format of event data as JSON format.
	//{
	//	"endTime":"2021-06-13T00:24:21.261+0800",
	//	"startTime":"2021-06-13T00:24:20.257+0800",
	//	"elapsedNano":1004326112,
	//	"timezone":"CST",
	//	"ids":{
	//	    "eventId":"72a59682-230f-4ba2-a9fc-e99a031e4d8c",
	//		"requestId":"",
	//		"traceId":""
	//  },
	//	"app":{
	//	    "appName":"appName",
	//		"appVersion":"unknown",
	//		"entryName":"unknown",
	//		"entryType":"unknown"
	//  },
	//	"env":{
	//	    "arch":"amd64",
	//		"hostname":"lark.local",
	//      "localIP":"localhost"
	//		"realm":"*",
	//		"region":"*",
	//		"az":"*",
	//		"domain":"*",
	//		"os":"darwin"
	//  },
	//	"payloads":{
	//	    "f1":"f2",
	//		"t2":"2021-06-13T00:24:21.261768+08:00"
	//  },
	//	"error":{
	//	    "my error":1
	//  },
	//	"counters":{
	//	    "count":1
	//  },
	//	"pairs":{
	//	    "key":"value"
	//  },
	//	"timing":{
	//	    "t1.count":1,
	//		"t1.elapsed_ms":1004
	//  },
	//	"remoteAddr":"localhost",
	//	"operation":"op",
	//	"eventStatus":"Ended",
	//	"resCode":"200"
	//}

	// endTime
	if event.GetEndTime().IsZero() {
		event.SetEndTime(time.Now())
	}
	// startTime
	if event.GetStartTime().IsZero() {
		event.SetStartTime(time.Now())
	}
	fields = append(fields,
		zap.Time(endTimeKey, event.GetEndTime()),
		zap.Time(startTimeKey, event.GetStartTime()),
		zap.Int64(elapsedKey, event.GetEndTime().Sub(event.GetStartTime()).Nanoseconds()),
		zap.String(timezoneKey, event.timeZone),
		zap.Any(idsKey, event.idsToMapObjectEncoder().Fields),
		zap.Any(appKey, event.appToMapObjectEncoder().Fields),
		zap.Any(envKey, event.envToMapObjectEncoder().Fields),
		zap.Any(payloadsKey, event.payloadsToMapObjectEncoder().Fields),
		zap.Any(errKey, event.errors.Fields),
		zap.Any(countersKey, event.counters.Fields),
		zap.Any(pairsKey, event.pairs.Fields),
		zap.Any(timingKey, event.timingToMapObjectEncoder().Fields),
		zap.String(remoteAddrKey, event.GetRemoteAddr()),
		zap.String(operationKey, event.GetOperation()),
		zap.String(eventStatusKey, event.GetEventStatus().String()))

	if len(event.errors.Fields) > 0 {
		fields = append(fields, zap.Any(errKey, event.errors.Fields))
	}

	// resCode
	if len(event.resCode) > 0 {
		fields = append(fields, zap.String(resCodeKey, event.resCode))
	}

	return fields
}

// Construct payloads to zapcore.MapObjectEncoder
func (event *eventZap) payloadsToMapObjectEncoder() *zapcore.MapObjectEncoder {
	enc := zapcore.NewMapObjectEncoder()
	for i := range event.payloads {
		event.payloads[i].AddTo(enc)
	}

	return enc
}

// Marshal payloads.
func (event *eventZap) marshalPayloads() string {
	return event.marshalEncoder(event.payloadsToMapObjectEncoder())
}

// Construct env to zapcore.MapObjectEncoder
func (event *eventZap) envToMapObjectEncoder() *zapcore.MapObjectEncoder {
	enc := zapcore.NewMapObjectEncoder()
	enc.AddString(hostnameKey, hostname)
	enc.AddString(localIpKey, localIp)
	enc.AddString(realmKey, realm)
	enc.AddString(regionKey, region)
	enc.AddString(azKey, az)
	enc.AddString(domainKey, domain)
	enc.AddString(goosKey, goos)
	enc.AddString(goArchKey, goArch)

	return enc
}

// Marshal env.
func (event *eventZap) marshalEnv() string {
	return event.marshalEncoder(event.envToMapObjectEncoder())
}

// Construct ids to zapcore.MapObjectEncoder
func (event *eventZap) idsToMapObjectEncoder() *zapcore.MapObjectEncoder {
	enc := zapcore.NewMapObjectEncoder()
	if len(event.eventId) > 0 {
		enc.AddString(eventIdKey, event.eventId)
	}

	if len(event.traceId) > 0 {
		enc.AddString(traceIdKey, event.traceId)
	}

	if len(event.requestId) > 0 {
		enc.AddString(requestIdKey, event.requestId)
	}

	return enc
}

// Marshal ids.
func (event *eventZap) marshalIds() string {
	return event.marshalEncoder(event.idsToMapObjectEncoder())
}

// Construct app to zapcore.MapObjectEncoder
func (event *eventZap) appToMapObjectEncoder() *zapcore.MapObjectEncoder {
	enc := zapcore.NewMapObjectEncoder()
	enc.AddString(appNameKey, event.appName)
	enc.AddString(appVersionKey, event.appVersion)
	enc.AddString(entryNameKey, event.entryName)
	enc.AddString(entryTypeKey, event.entryType)

	return enc
}

// Marshal app.
func (event *eventZap) marshalApp() string {
	return event.marshalEncoder(event.appToMapObjectEncoder())
}

// Construct timing to zapcore.MapObjectEncoder
func (event *eventZap) timingToMapObjectEncoder() *zapcore.MapObjectEncoder {
	enc := zapcore.NewMapObjectEncoder()
	for _, v := range event.tracker {
		v.ToZapFields(enc)
	}

	return enc
}

// Marshal timing.
func (event *eventZap) marshalTiming() string {
	return event.marshalEncoder(event.timingToMapObjectEncoder())
}

// Marshal zapcore.MapObjectEncoder.
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

// Is Event in progress?
func (event *eventZap) inProgress() bool {
	if event.status != InProgress {
		return false
	}

	return true
}

// Convert time.Time to milliseconds.
func toMillisecond(curr time.Time) int64 {
	return curr.UnixNano() / 1e6
}
