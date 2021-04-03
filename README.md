<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [rk-query](#rk-query)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
    - [With Rk format](#with-rk-format)
    - [With JSON format](#with-json-format)
    - [Development Status: Stable](#development-status-stable)
    - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# rk-query
Human readable query logger with zap logger

- [zap](https://github.com/uber-go/zap)
- [lumberjack](https://github.com/natefinch/lumberjack)
- [rk-logger](https://github.com/rookie-ninja/rk-logger)

## Installation
`go get -u rookie-ninja/rk-query`

## Quick Start
Zap logger needs to be pass to query in order to write logs

Please refer https://github.com/rookie-ninja/rk-logger for easy initialization of zap logger

### With Rk format
It is human friendly printed query log format

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
    "maxsize": 1,
    "maxage": 7,
    "maxbackups": 3,
    "localtime": true,
    "compress": true
   }`)
)

func withEventRkFormat() {
	logger, _, _ := rk_logger.NewZapLoggerWithBytes(bytes, rk_logger.JSON)

	fac := rk_query.NewEventFactory(
		rk_query.WithAppName("appName"),
		rk_query.WithFormat(rk_query.RK),
		rk_query.WithOperation("op"),
		rk_query.WithLogger(logger))
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
end_time=2020-07-30T03:42:22.393874+08:00
start_time=2020-07-30T03:42:21.390642+08:00
elapsed_nano=1004000000
hostname=MYLOCAL
timing={"t1.count":1,"t1.elapsed_ms":1003}
counter={"count":1}
pair={"key":"value"}
error={"MyError":1}
field={"f1":"f2","t2":"2020-07-30T03:42:22.393857+08:00"}
remote_addr=Unknown
app_name=appName
operation=Unknown
event_status=Ended
history=s-t1:1596051741390,e-t1:1003,end:0
EOE
```

### With JSON format
It is parsing friendly printed query log format

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
    "maxsize": 1,
    "maxage": 7,
    "maxbackups": 3,
    "localtime": true,
    "compress": true
   }`)
)

func withEventJSONFormat() {
	logger, _, _ := rk_logger.NewZapLoggerWithBytes(bytes, rk_logger.JSON)

	fac := rk_query.NewEventFactory(
		rk_query.WithAppName("appName"),
		rk_query.WithFormat(rk_query.JSON),
		rk_query.WithOperation("op"),
		rk_query.WithLogger(logger))
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
    "end_time":"2020-07-30T03:42:23.398+0800",
    "start_time":"2020-07-30T03:42:22.394+0800",
    "elapsed_nano":1004000000,
    "hostname":"MYLOCAL",
    "timing":{
        "t1.count":1,
        "t1.elapsed_ms":1004
    },
    "counter":{
        "count":1
    },
    "pair":{
        "key":"value"
    },
    "error":{
        "MyError":1
    },
    "field":{
        "f1":"f2",
        "t2":"2020-07-30T03:42:23.398282+08:00"
    },
    "remote_addr":"Unknown",
    "app_name":"appName",
    "operation":"op",
    "event_status":"Ended",
    "history":"s-t1:1596051742394,e-t1:1004,end:0"
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

<hr>

Released under the [MIT License](LICENSE).

