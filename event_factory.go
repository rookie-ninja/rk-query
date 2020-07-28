// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"go.uber.org/zap"
	"os"
)

// Not thread safe!!!
type EventFactory struct {
	TimeSource        TimeSource
	AppName           string
	HostName          string
	Format            Format
	Minimal           bool
	ZapLogger         *zap.Logger
	Listeners         []eventEntryListener
	DefaultKvs        map[string]string
}

func NewEventFactory(app string, ts TimeSource, logger *zap.Logger) *EventFactory {
	return &EventFactory{
		TimeSource:        ts,
		AppName:           app,
		HostName:          obtainHostName(),
		ZapLogger:         logger,
		Format:            RK,
		Minimal:           false,
		Listeners:         make([]eventEntryListener, 0),
		DefaultKvs: make(map[string]string),
	}
}

func (factory *EventFactory) CreateEvent() Event {
	event := NewEventImpl(
		factory.TimeSource,
		factory.AppName,
		factory.HostName,
		factory.DefaultKvs,
		factory.ZapLogger,
		factory.Format,
		factory.Minimal,
		factory.Listeners,
		false)

	return event
}

func (factory *EventFactory) CreateThreadSafeEvent() Event {
	return NewThreadSafeEventImpl(factory.CreateEvent())
}

func (factory *EventFactory) CreateNoopEvent() Event {
	event := NoopEvent{}
	return &event
}

func obtainHostName() string {
	hostName, err := os.Hostname()

	// In this version, we will ignore errors returned by OS
	if err != nil {
		hostName = Unknown
	}

	return hostName
}
