// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"github.com/rookie-ninja/rk-logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

type EventOption func(Event)

func WithLogger(logger *zap.Logger) EventOption {
	return func(event Event) {
		event.setLogger(logger)
	}
}

func WithFormat(format format) EventOption {
	return func(event Event) {
		event.setFormat(format)
	}
}

func WithQuietMode(quietMode bool) EventOption {
	return func(event Event) {
		event.setQuietMode(quietMode)
	}
}

func WithAppName(appName string) EventOption {
	return func(event Event) {
		event.setAppName(appName)
	}
}

func WithHostname(hostname string) EventOption {
	return func(event Event) {
		event.setHostname(hostname)
	}
}

func WithOperation(operation string) EventOption {
	return func(event Event) {
		event.SetOperation(operation)
	}
}

func WithRemoteAddr(addr string) EventOption {
	return func(event Event) {
		event.SetRemoteAddr(addr)
	}
}

func WithFields(fields []zap.Field) EventOption {
	return func(event Event) {
		event.AddFields(fields...)
	}
}

// Not thread safe!!!
type EventFactory struct {
	appName string
	options []EventOption
}

func NewEventFactory(option ...EventOption) *EventFactory {
	factory := &EventFactory{
		options: option,
	}

	return factory
}

func (factory *EventFactory) GetAppName() string {
	return factory.appName
}

func (factory *EventFactory) CreateEvent(options ...EventOption) Event {
	event := &eventZap{
		logger:     rklogger.EventLogger,
		format:     RK,
		status:     notStarted,
		appName:    unknown,
		hostname:   obtainHostName(),
		remoteAddr: obtainHostName(),
		operation:  unknown,
		counters:   zapcore.NewMapObjectEncoder(),
		pairs:      zapcore.NewMapObjectEncoder(),
		errors:     zapcore.NewMapObjectEncoder(),
		fields:     make([]zap.Field, 0),
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

	factory.appName = event.GetAppName()

	event.logger.Core().Sync()

	if !event.quietMode {
		event.eventHistory = newEventHistory()
	}

	return event
}

func (factory *EventFactory) CreateEventNoop() Event {
	return &eventNoop{}
}

func (factory *EventFactory) CreateEventThreadSafe(options ...EventOption) Event {
	event := factory.CreateEvent(options...)
	return &eventThreadSafe{
		delegate: event,
		lock:     &sync.Mutex{},
	}
}

func obtainHostName() string {
	hostName, err := os.Hostname()

	// In this version, we will ignore errors returned by OS
	if err != nil {
		hostName = unknown
	}

	return hostName
}
