package main

import (
	rk_logger "github.com/rookie-ninja/rk-logger"
	rk_query "github.com/rookie-ninja/rk-query"
	"time"
)

func main() {
	withRawEventRkFormat()
	withRawEventRkMinFormat()
	withEventHelper()
}

func withEventHelper() {
	zapBytes := []byte(`{
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
      }
    }`)

	lumberBytes := []byte(`{
     "maxsize": 1,
     "maxage": 7,
     "maxbackups": 3,
     "localtime": true,
     "compress": true
    }`)

	lumber, err := rk_logger.NewLumberjackLoggerWithBytes(lumberBytes, rk_logger.JSON)
	if err != nil {
		panic(err)
	}

	logger, _, err := rk_logger.NewZapLoggerWithBytes(zapBytes, rk_logger.JSON, lumber)

	helper := rk_query.NewEventHelperWithLogger("fake", logger)
	event := helper.Start("my-op")
	event.SetRemoteAddr("1.1.1.1")
	event.StartTimer("t1")
	time.Sleep(1 * time.Second)
	event.EndTimer("t1")

	event.AddErr(MyError{})
	helper.Finish(event)
}

func withRawEventRkFormat() {
	zapBytes := []byte(`{
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
      }
    }`)

	lumberBytes := []byte(`{
     "maxsize": 1,
     "maxage": 7,
     "maxbackups": 3,
     "localtime": true,
     "compress": true
    }`)

	lumber, err := rk_logger.NewLumberjackLoggerWithBytes(lumberBytes, rk_logger.JSON)
	if err != nil {
		panic(err)
	}

	logger, _, err := rk_logger.NewZapLoggerWithBytes(zapBytes, rk_logger.JSON, lumber)

	ts := rk_query.RealTimeSource{}

	fac := rk_query.EventFactory{
		TimeSource: &ts,
		AppName:    "my-app",
		Format:     rk_query.RK,
		Minimal:    false,
		ZapLogger:  logger,
		HostName:   "my-host",
	}

	event := fac.CreateEvent()
	event.SetRemoteAddr("1.1.1.1")
	event.SetOperation("my-op")
	event.SetStartTimeMS(ts.CurrentTimeMS())
	event.StartTimer("t1")
	time.Sleep(1 * time.Second)
	event.EndTimer("t1")

	event.AddErr(MyError{})
	event.SetEndTimeMS(ts.CurrentTimeMS())
	event.WriteLog()
}

func withRawEventRkMinFormat() {
	zapBytes := []byte(`{
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
      }
    }`)

	lumberBytes := []byte(`{
     "maxsize": 1,
     "maxage": 7,
     "maxbackups": 3,
     "localtime": true,
     "compress": true
    }`)

	lumber, err := rk_logger.NewLumberjackLoggerWithBytes(lumberBytes, rk_logger.JSON)
	if err != nil {
		panic(err)
	}

	logger, _, err := rk_logger.NewZapLoggerWithBytes(zapBytes, rk_logger.JSON, lumber)

	ts := rk_query.RealTimeSource{}

	fac := rk_query.EventFactory{
		TimeSource: &ts,
		AppName:    "my-app",
		Format:     rk_query.RK,
		Minimal:    true,
		ZapLogger:  logger,
		HostName:   "my-host",
	}

	event := fac.CreateEvent()
	event.SetRemoteAddr("1.1.1.1")
	event.SetOperation("my-op")
	event.SetStartTimeMS(ts.CurrentTimeMS())
	event.StartTimer("t1")
	time.Sleep(1 * time.Second)
	event.EndTimer("t1")

	event.AddErr(MyError{})
	event.SetEndTimeMS(ts.CurrentTimeMS())
	event.WriteLog()
}

type MyError struct {
}

func (err MyError) Error() string {
	return ""
}
