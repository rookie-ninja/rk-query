// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkquery

import (
	"github.com/google/uuid"
	"github.com/rookie-ninja/rk-logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	goos   = runtime.GOOS
	goArch = runtime.GOARCH
)

var (
	domain   = ""
	localIp  = getLocalIP()
	hostname = getHostName()
)

// EventOption will be pass into EventFactory while creating Event to override fields in Event.
type EventOption func(Event)

// WithZapLogger override logger in Event.
func WithZapLogger(logger *zap.Logger) EventOption {
	return func(event Event) {
		if logger == nil {
			return
		}

		switch v := event.(type) {
		case *eventZap:
			v.logger = logger
		case *eventThreadSafe:
			v.delegate.logger = logger
		}
	}
}

// WithEncoding override encoding in Event.
func WithEncoding(ec Encoding) EventOption {
	return func(event Event) {
		if ec != JSON && ec != CONSOLE && ec != FLATTEN {
			return
		}

		switch v := event.(type) {
		case *eventZap:
			v.encoding = ec
		case *eventThreadSafe:
			v.delegate.encoding = ec
		}
	}
}

// WithQuietMode turn on quiet mode which won't flush data to logger.
func WithQuietMode(quietMode bool) EventOption {
	return func(event Event) {
		switch v := event.(type) {
		case *eventZap:
			v.quietMode = quietMode
		case *eventThreadSafe:
			v.delegate.quietMode = quietMode
		}
	}
}

// WithEntryName override entry name in Event.
func WithEntryName(entryName string) EventOption {
	return func(event Event) {
		switch v := event.(type) {
		case *eventZap:
			v.entryName = entryName
		case *eventThreadSafe:
			v.delegate.entryName = entryName
		}
	}
}

// WithEntryType override entry type in Event.
func WithEntryType(entryType string) EventOption {
	return func(event Event) {
		switch v := event.(type) {
		case *eventZap:
			v.entryType = entryType
		case *eventThreadSafe:
			v.delegate.entryType = entryType
		}
	}
}

// WithAppName override app name in Event.
func WithAppName(appName string) EventOption {
	return func(event Event) {
		switch v := event.(type) {
		case *eventZap:
			v.appName = appName
		case *eventThreadSafe:
			v.delegate.appName = appName
		}
	}
}

// WithAppVersion overrides app version in event.
func WithAppVersion(appVersion string) EventOption {
	return func(event Event) {
		switch v := event.(type) {
		case *eventZap:
			v.appVersion = appVersion
		case *eventThreadSafe:
			v.delegate.appVersion = appVersion
		}
	}
}

// WithOperation overrides operation in Event.
func WithOperation(operation string) EventOption {
	return func(event Event) {
		event.SetOperation(operation)
	}
}

// WithPayloads overrides payloads with form of zap.Field in Event.
func WithPayloads(fields ...zap.Field) EventOption {
	return func(event Event) {
		event.AddPayloads(fields...)
	}
}

// EventFactory is not thread safe!!!
type EventFactory struct {
	options []EventOption
}

// NewEventFactory creates a new event factory with option.
func NewEventFactory(option ...EventOption) *EventFactory {
	factory := &EventFactory{
		options: option,
	}

	domain = getDefaultIfEmptyString(os.Getenv("DOMAIN"), "*")

	return factory
}

// CreateEvent creates a new event with options.
func (factory *EventFactory) CreateEvent(options ...EventOption) Event {
	event := &eventZap{
		logger:     rklogger.EventLogger,
		encoding:   CONSOLE,
		quietMode:  false,
		appName:    "",
		appVersion: "",
		entryName:  "",
		entryType:  "",
		eventId:    generateEventId(),
		traceId:    "",
		requestId:  "",
		startTime:  time.Now(),
		timeZone:   getTimeZone(),
		payloads:   make([]zap.Field, 0),
		errors:     zapcore.NewMapObjectEncoder(),
		operation:  "",
		remoteAddr: "localhost",
		resCode:    "",
		status:     NotStarted,
		pairs:      zapcore.NewMapObjectEncoder(),
		counters:   zapcore.NewMapObjectEncoder(),
		tracker:    make(map[string]*timeTracker),
	}

	for i := range factory.options {
		opt := factory.options[i]
		opt(event)
	}

	for i := range options {
		opt := options[i]
		opt(event)
	}

	event.logger.Core().Sync()

	return event
}

// CreateEventNoop creates a new noop event.
func (factory *EventFactory) CreateEventNoop() Event {
	return &eventNoop{}
}

// CreateEventThreadSafe create a new thread safe event.
func (factory *EventFactory) CreateEventThreadSafe(options ...EventOption) Event {
	event := factory.CreateEvent(options...)
	return &eventThreadSafe{
		delegate: event.(*eventZap),
		lock:     &sync.Mutex{},
	}
}

// Get hostname of current machine.
func getHostName() string {
	hostName, err := os.Hostname()

	// In this version, we will ignore errors returned by OS
	if err != nil {
		hostName = unknown
	}

	return hostName
}

// Generate request id based on google/uuid.
// UUIDs are based on RFC 4122 and DCE 1.1: Authentication and Security Services.
//
// A UUID is a 16 byte (128 bit) array. UUIDs may be used as keys to maps or compared directly.
func generateEventId() string {
	// do not use uuid.New() since it would panic if any error occurs
	requestId, err := uuid.NewRandom()

	// currently, we will return empty string if error occurs
	if err != nil {
		return ""
	}

	return requestId.String()
}

// Get time zone.
func getTimeZone() string {
	zone, _ := time.Now().Zone()
	return zone
}

// Return default value if original string is empty.
func getDefaultIfEmptyString(origin, def string) string {
	if len(origin) < 1 {
		return def
	}

	return origin
}

// This is a tricky function.
// We will iterate through all the network interfaces，but will choose the first one since we are assuming that
// eth0 will be the default one to use in most of the case.
//
// Currently, we do not have any interfaces for selecting the network interface yet.
func getLocalIP() string {
	localIP := "localhost"

	// skip the error since we don't want to break RPC calls because of it
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return localIP
	}

	for _, addr := range addresses {
		items := strings.Split(addr.String(), "/")
		if len(items) < 2 || items[0] == "127.0.0.1" {
			continue
		}

		if match, err := regexp.MatchString(`\d+\.\d+\.\d+\.\d+`, items[0]); err == nil && match {
			localIP = items[0]
		}
	}

	return localIP
}
