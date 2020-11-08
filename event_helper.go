// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	rk_logger "github.com/rookie-ninja/rk-logger"
	"time"
)

var (
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

	StdoutLogger, _, _ = rk_logger.NewZapLoggerWithBytes(StdLoggerConfigBytes, rk_logger.JSON)
)

// A helper function for easy use of EventData
type eventHelper struct {
	Factory *EventFactory
}

func NewEventHelper(factory *EventFactory) *eventHelper {
	if factory == nil {
		factory = NewEventFactory()
	}
	return &eventHelper{factory}
}

func (helper *eventHelper) Start(operation string) Event {
	event := helper.Factory.CreateEvent()

	event.SetOperation(operation)
	event.SetStartTime(time.Now())
	return event
}

func (helper *eventHelper) Finish(event Event) {
	event.SetEndTime(time.Now())
	event.WriteLog()
}

func (helper *eventHelper) FinishWithCond(event Event, success bool) {
	if success {
		event.SetCounter("success", 1)
	} else {
		event.SetCounter("failure", 1)
	}

	helper.Finish(event)
}

func (helper *eventHelper) FinishWithError(event Event, err error) {
	if err == nil {
		helper.FinishWithCond(event, true)
	}

	event.AddErr(err)
	helper.FinishWithCond(event, false)
}
