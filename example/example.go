// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
	"time"
)

var (
	bytes = []byte(`{
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
    "maxsize": 1024,
    "maxage": 7,
    "maxbackups": 3,
    "localtime": true,
    "compress": true
   }`)
)

func main() {
	withEventConsoleEncoding()
	withEventJSONEncoding()
	withEventHelper()
}

func withEventJSONEncoding() {
	logger, _, _ := rklogger.NewZapLoggerWithBytes(bytes, rklogger.JSON)

	fac := rkquery.NewEventFactory(
		rkquery.WithAppName("appName"),
		rkquery.WithEncoding(rkquery.JSON),
		rkquery.WithOperation("op"),
		rkquery.WithZapLogger(logger))
	event := fac.CreateEvent()

	event.SetStartTime(time.Now())
	event.StartTimer("t1")
	time.Sleep(1 * time.Second)
	event.EndTimer("t1")
	event.AddPair("key", "value")
	event.SetCounter("count", 1)
	event.AddPayloads(zap.String("f1", "f2"), zap.Time("t2", time.Now()))
	event.AddErr(MyError{})
	event.SetResCode("200")
	event.SetEndTime(time.Now())
	event.Finish()
}

func withEventConsoleEncoding() {
	logger, _, _ := rklogger.NewZapLoggerWithBytes(bytes, rklogger.JSON)

	fac := rkquery.NewEventFactory(
		rkquery.WithEntryName("entry-example"),
		rkquery.WithEntryType("example"),
		rkquery.WithAppName("appName"),
		rkquery.WithAppVersion("v0.0.1"),
		rkquery.WithEncoding(rkquery.CONSOLE),
		rkquery.WithOperation("op"),
		rkquery.WithZapLogger(logger))
	event := fac.CreateEvent()

	event.SetStartTime(time.Now())
	event.StartTimer("t1")
	time.Sleep(1 * time.Second)
	event.EndTimer("t1")
	event.AddPair("key", "value")
	event.SetCounter("count", 1)
	event.AddPayloads(zap.String("f1", "f2"), zap.Time("t2", time.Now()))
	event.AddErr(MyError{})
	event.SetResCode("200")
	event.SetEndTime(time.Now())
	event.Finish()
}

func withEventHelper() {
	logger, _, _ := rklogger.NewZapLoggerWithBytes(bytes, rklogger.JSON)
	helper := rkquery.NewEventHelper(rkquery.NewEventFactory(rkquery.WithZapLogger(logger)))

	event := helper.Start("op")
	helper.Finish(event)
}

type MyError struct{}

func (err MyError) Error() string {
	return "my error"
}
