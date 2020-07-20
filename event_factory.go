package rk_query

import (
	"github.com/juju/errors"
	"go.uber.org/zap"
	"os"
)

// Not thread safe!!!
type EventFactory struct {
	TimeSource        TimeSource
	AppName           string
	HostName          string
	Format            Format
	Minimal			  bool
	ZapLogger         *zap.Logger
	Listeners         []eventEntryListener
	DefaultNameValues map[string]string
}

func NewEventFactory(logger *zap.Logger) (*EventFactory, error) {
	if logger == nil {
		return nil, errors.NewNotValid(nil, "zap logger is nil")
	}

	return &EventFactory{
		TimeSource:        &RealTimeSource{},
		AppName:   		   EventUnknownApplication,
		HostName:          obtainHostName(),
		ZapLogger:         logger,
		Format: 		   RK,
		Minimal:           false,
		Listeners:         make([]eventEntryListener, 0),
		DefaultNameValues: make(map[string]string),
	}, nil
}

func (factory *EventFactory) CreateEvent() Event {
	event := NewEventImpl(
		factory.TimeSource,
		factory.AppName,
		factory.HostName,
		factory.DefaultNameValues,
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
	event := NoopEventData{}
	return &event
}

func obtainHostName() string {
	hostName, err := os.Hostname()

	// In this version, we will ignore errors returned by OS
	if err != nil {
		hostName = EventUnknownHostName
	}

	return hostName
}