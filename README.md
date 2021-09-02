# rk-query
[![build](https://github.com/rookie-ninja/rk-query/actions/workflows/ci.yml/badge.svg)](https://github.com/rookie-ninja/rk-query/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Human readable query logger with [zap](https://github.com/uber-go/zap), [lumberjack](https://github.com/natefinch/lumberjack) and [rk-logger](https://github.com/rookie-ninja/rk-logger)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
- [Quick Start](#quick-start)
  - [With Console encoding](#with-console-encoding)
  - [With JSON encoding](#with-json-encoding)
  - [Development Status: Stable](#development-status-stable)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
`go get -u github.com/rookie-ninja/rk-query`

## Quick Start
Zap logger needs to be pass to query in order to write logs

Please refer https://github.com/rookie-ninja/rk-logger for easy initialization of zap logger

### With Console encoding
It is human friendly printed query log encoding type.

Example:
```go
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

func withEventConsoleEncoding() {
	logger, _, _ := rklogger.NewZapLoggerWithBytes(bytes, rk_logger.JSON)

	fac := rkquery.NewEventFactory(
		rkquery.WithAppName("appName"),
		rkquery.WithEncoding(rkquery.CONSOLE),
		rkquery.WithOperation("op"),
		rkquery.WithLogger(logger))
	event := fac.CreateEvent()

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
```
Output
```
------------------------------------------------------------------------
endTime=2021-06-13T01:16:27.58556+08:00
startTime=2021-06-13T01:16:26.581691+08:00
elapsedNano=1003868481
timezone=CST
ids={"eventId":"581812ae-924a-44b2-83f8-fa8eef071393"}
app={"appName":"appName","appVersion":"v0.0.1","entryName":"entry-example","entryType":"example"}
env={"arch":"amd64","hostname":"lark.local","realm":"*","region":"*","az":"*","domain":"*","os":"darwin"}
payloads={"f1":"f2","t2":"2021-06-13T01:16:27.58554+08:00"}
error={"my error":1}
counters={"count":1}
pairs={"key":"value"}
timing={"t1.count":1,"t1.elapsedMs":1004}
remoteAddr=localhost
operation=op
resCode=200
eventStatus=Ended
EOE
```

### With JSON encoding
It is parsing friendly printed query log encoding type.

Example:
```go
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

func withEventJSONEncoding() {
	logger, _, _ := rklogger.NewZapLoggerWithBytes(bytes, rk_logger.JSON)

	fac := rkquery.NewEventFactory(
		rkquery.WithAppName("appName"),
		rkquery.WithEncoding(rkquery.JSON),
		rkquery.WithOperation("op"),
		rkquery.WithLogger(logger))
	event := fac.CreateEvent()

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
```
Output 
We formatted JSON output bellow, actual logs would not be a pretty formatted JSON
```
{
    "endTime":"2021-06-13T00:24:21.261+0800",
    "startTime":"2021-06-13T00:24:20.257+0800",
    "elapsedNano":1004326112,
    "timezone":"CST",
    "ids":{
        "eventId":"72a59682-230f-4ba2-a9fc-e99a031e4d8c",
        "requestId":"",
        "traceId":""
    },
    "app":{
        "appName":"appName",
        "appVersion":"unknown",
        "entryName":"unknown",
        "entryType":"unknown"
    },
    "env":{
        "arch":"amd64",
        "hostname":"lark.local",
        "realm":"*",
        "region":"*",
        "az":"*",
        "domain":"*",
        "os":"darwin"
    },
    "payloads":{
        "f1":"f2",
        "t2":"2021-06-13T00:24:21.261768+08:00"
    },
    "error":{
        "my error":1
    },
    "counters":{
        "count":1
    },
    "pairs":{
        "key":"value"
    },
    "timing":{
        "t1.count":1,
        "t1.elapsed_ms":1004
    },
    "remoteAddr":"localhost",
    "operation":"op",
    "eventStatus":"Ended",
    "resCode":"200"
}
```

### Development Status: Stable

### Contributing
We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The rk maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
dongxuny@gmail.com. That email list is a private, safe space; even the zap
maintainers don't have access, so don't hesitate to hold us to a high
standard.

Released under the [Apache 2.0 License](LICENSE).

