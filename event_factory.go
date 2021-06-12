// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"github.com/google/uuid"
	"github.com/rookie-ninja/rk-logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	goos   = runtime.GOOS
	goArch = runtime.GOARCH
)

type EventOption func(Event)

// Provide zap logger.
func WithZapLogger(logger *zap.Logger) EventOption {
	return func(event Event) {
		switch v := event.(type) {
		case *eventZap:
			v.logger = logger
		case *eventThreadSafe:
			v.delegate.logger = logger
		}
	}
}

// Provide format.
func WithFormat(format format) EventOption {
	return func(event Event) {
		switch v := event.(type) {
		case *eventZap:
			v.format = format
		case *eventThreadSafe:
			v.delegate.format = format
		}
	}
}

// Turn on quiet mode which won't flush data to logger.
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

// Provide entry name.
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

// Provide entry type.
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

// Provide app name.
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

// Provide app version.
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

// Provide locale.
func WithLocale(locale string) EventOption {
	return func(event Event) {
		switch v := event.(type) {
		case *eventZap:
			v.locale = locale
		case *eventThreadSafe:
			v.delegate.locale = locale
		}
	}
}

// Provide operation.
func WithOperation(operation string) EventOption {
	return func(event Event) {
		event.SetOperation(operation)
	}
}

// Provide payloads with form of zap.Field.
func WithPayloads(fields ...zap.Field) EventOption {
	return func(event Event) {
		event.AddPayloads(fields...)
	}
}

// Not thread safe!!!
type EventFactory struct {
	options []EventOption
}

// Create a new event factory with option.
func NewEventFactory(option ...EventOption) *EventFactory {
	factory := &EventFactory{
		options: option,
	}

	return factory
}

// Create a new event with option.
func (factory *EventFactory) CreateEvent(options ...EventOption) Event {
	event := &eventZap{
		logger:     rklogger.EventLogger,
		format:     RK,
		quietMode:  false,
		appName:    unknown,
		appVersion: unknown,
		entryName:  unknown,
		entryType:  unknown,
		hostname:   obtainHostName(),
		locale:     getLocale(),
		eventId:    generateEventId(),
		traceId:    "",
		requestId:  "",
		startTime:  time.Now(),
		timeZone:   getTimeZone(),
		payloads:   make([]zap.Field, 0),
		errors:     zapcore.NewMapObjectEncoder(),
		operation:  unknown,
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

// Create a new noop event.
func (factory *EventFactory) CreateEventNoop() Event {
	return &eventNoop{}
}

// Create a new thread safe event.
func (factory *EventFactory) CreateEventThreadSafe(options ...EventOption) Event {
	event := factory.CreateEvent(options...)
	return &eventThreadSafe{
		delegate: event.(*eventZap),
		lock:     &sync.Mutex{},
	}
}

// Get hostname of current machine.
func obtainHostName() string {
	hostName, err := os.Hostname()

	// In this version, we will ignore errors returned by OS
	if err != nil {
		hostName = unknown
	}

	return hostName
}

// Get locale from environment variable
func getLocale() string {
	realm, region, az, domain := "*", "*", "*", "*"
	if v := os.Getenv("REALM"); len(v) > 0 {
		realm = v
	}
	if v := os.Getenv("REGION"); len(v) > 0 {
		region = v
	}
	if v := os.Getenv("AZ"); len(v) > 0 {
		az = v
	}
	if v := os.Getenv("DOMAIN"); len(v) > 0 {
		domain = v
	}

	elements := []string{realm, region, az, domain}

	return strings.Join(elements, "::")
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
