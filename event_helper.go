// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkquery

import (
	rk_logger "github.com/rookie-ninja/rk-logger"
	"time"
)

var (
	// StdLoggerConfigBytes defines zap logger config whose output path is stdout.
	StdLoggerConfigBytes = []byte(`{
     "level": "info",
     "encoding": "console",
     "outputPaths": ["stdout"],
     "errorOutputPaths": ["stderr"],
     "initialFields": {},
     "encoderConfig": {
       "messageKey": "msg",
       "levelKey": "",
       "nameKey": "",
       "timeKey": "",
       "callerKey": "",
       "stacktraceKey": "",
       "callstackKey": "",
       "errorKey": "",
       "timeEncoder": "iso8601",
       "fileKey": "",
       "levelEncoder": "capital",
       "durationEncoder": "second",
       "callerEncoder": "full",
       "nameEncoder": "full"
     },
    "maxsize": 1,
    "maxage": 7,
    "maxbackups": 3,
    "localtime": true,
    "compress": true
   }`)

	// StdoutLogger defines zap logger which use StdLoggerConfigBytes as config.
	StdoutLogger, _, _ = rk_logger.NewZapLoggerWithBytes(StdLoggerConfigBytes, rk_logger.JSON)
)

// EventHelper is a helper function for easy use of EventData.
type EventHelper struct {
	Factory *EventFactory
}

// NewEventHelper creates a new event helper.
func NewEventHelper(factory *EventFactory) *EventHelper {
	if factory == nil {
		factory = NewEventFactory()
	}
	return &EventHelper{factory}
}

// Start function creates and start a new event with options.
func (helper *EventHelper) Start(operation string, opts ...EventOption) Event {
	event := helper.Factory.CreateEvent(opts...)

	event.SetOperation(operation)
	event.SetStartTime(time.Now())
	return event
}

// Finish current event.
func (helper *EventHelper) Finish(event Event) {
	event.SetResCode("OK")
	event.SetEndTime(time.Now())
	event.Finish()
}

// FinishWithCond finish current event with condition.
func (helper *EventHelper) FinishWithCond(event Event, success bool) {
	if success {
		event.SetCounter("success", 1)
		event.SetResCode("OK")
	} else {
		event.SetCounter("failure", 1)
		event.SetResCode("Fail")
	}

	event.SetEndTime(time.Now())
	event.Finish()
}

// FinishWithError finish current event with error.
func (helper *EventHelper) FinishWithError(event Event, err error) {
	if err == nil {
		helper.FinishWithCond(event, true)
	}

	event.SetResCode("Fail")

	event.AddErr(err)
	helper.FinishWithCond(event, false)
}
