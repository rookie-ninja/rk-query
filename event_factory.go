// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type EventZapOption func(*EventZap)

func WithLogger(logger *zap.Logger) EventZapOption {
	return func(event *EventZap) {
		event.logger = logger
	}
}

func WithFormat(format Format) EventZapOption {
	return func(event *EventZap) {
		event.format = format
	}
}

func WithQuietMode(quietMode bool) EventZapOption {
	return func(event *EventZap) {
		event.quietMode = quietMode
	}
}

func WithAppName(appName string) EventZapOption {
	return func(event *EventZap) {
		event.appName = appName
	}
}

func WithHostname(hostname string) EventZapOption {
	return func(event *EventZap) {
		event.hostname = hostname
	}
}

func WithOperation(operation string) EventZapOption {
	return func(event *EventZap) {
		event.operation = operation
	}
}

func WithRemoteAddr(addr string) EventZapOption {
	return func(event *EventZap) {
		event.remoteAddr = addr
	}
}

func WithFields(fields []zap.Field) EventZapOption {
	return func(event *EventZap) {
		event.fields = append(event.fields, fields...)
	}
}

// Not thread safe!!!
type EventZapFactory struct {
	options []EventZapOption
}

func NewEventZapFactory(option ...EventZapOption) *EventZapFactory {
	return &EventZapFactory{
		options: option,
	}
}

func (factory *EventZapFactory) CreateEventZap(options ...EventZapOption) *EventZap {
	event := &EventZap{
		logger:     zap.NewNop(),
		format:     RK,
		status:     notStarted,
		appName:    Unknown,
		hostname:   obtainHostName(),
		remoteAddr: Unknown,
		operation:  Unknown,
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
		opt := factory.options[i]
		opt(event)
	}

	event.logger.Core().Sync()

	if !event.quietMode {
		event.eventHistory = newEventHistory()
	}

	return event
}

func obtainHostName() string {
	hostName, err := os.Hostname()

	// In this version, we will ignore errors returned by OS
	if err != nil {
		hostName = Unknown
	}

	return hostName
}
