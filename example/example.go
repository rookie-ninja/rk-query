// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
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
    "maxsize": 1,
    "maxage": 7,
    "maxbackups": 3,
    "localtime": true,
    "compress": true
   }`)
)

func main() {
	withEventZapRkFormat()
	withEventZapJSONFormat()
	withEventZapHelper()
}

func withEventZapJSONFormat() {
	logger, _, _ := rk_logger.NewZapLoggerWithBytes(bytes, rk_logger.JSON)

	fac := rk_query.NewEventZapFactory(
		rk_query.WithAppName("appName"),
		rk_query.WithFormat(rk_query.JSON),
		rk_query.WithOperation("op"),
		rk_query.WithLogger(logger))
	event := fac.CreateEventZap()

	event.SetStartTime(time.Now())
	event.StartTimer("t1")
	time.Sleep(1 * time.Second)
	event.EndTimer("t1")
	event.AddPair("key", "value")
	event.SetCounter("count", 1)
	event.AddFields(zap.String("f1", "f2"), zap.Time("t2", time.Now()))
	event.AddErr(MyError{})
	event.SetEndTime(time.Now())
	event.WriteLog()
}

func withEventZapRkFormat() {
	logger, _, _ := rk_logger.NewZapLoggerWithBytes(bytes, rk_logger.JSON)

	fac := rk_query.NewEventZapFactory(
		rk_query.WithAppName("appName"),
		rk_query.WithFormat(rk_query.RK),
		rk_query.WithOperation("op"),
		rk_query.WithLogger(logger))
	event := fac.CreateEventZap()

	event.SetStartTime(time.Now())
	event.StartTimer("t1")
	time.Sleep(1 * time.Second)
	event.EndTimer("t1")
	event.AddPair("key", "value")
	event.SetCounter("count", 1)
	event.AddFields(zap.String("f1", "f2"), zap.Time("t2", time.Now()))
	event.AddErr(MyError{})
	event.SetEndTime(time.Now())
	event.WriteLog()
}

func withEventZapHelper() {
	logger, _, _ := rk_logger.NewZapLoggerWithBytes(bytes, rk_logger.JSON)
	helper := rk_query.NewEventZapHelper(rk_query.NewEventZapFactory(rk_query.WithLogger(logger)))

	event := helper.Start("op")
	helper.Finish(event)
}

type MyError struct {
}

func (err MyError) Error() string {
	return ""
}
