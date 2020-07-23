// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import "go.uber.org/zap"

// A helper function for easy use of EventData
type EventHelper struct {
	Factory    *EventFactory
	TimeSource TimeSource
}

func NewEventHelperWithZapLogger(appName string, ts TimeSource, logger *zap.Logger) *EventHelper {
	factory := NewEventFactory(appName, ts, logger)

	return &EventHelper{factory, ts}
}

func (helper *EventHelper) Start(operationName string) Event {
	event := helper.Factory.CreateEvent()

	event.SetOperation(operationName)
	event.SetStartTimeMS(helper.TimeSource.CurrentTimeMS())
	return event
}

func (helper *EventHelper) Finish(event Event) {
	event.SetEndTimeMS(helper.TimeSource.CurrentTimeMS())
	event.WriteLog()
}

func (helper *EventHelper) FinishWithCond(event Event, success bool) {
	if success {
		event.SetCounter("success", 1)
	} else {
		event.SetCounter("failure", 1)
	}

	helper.Finish(event)
}

func (helper *EventHelper) FinishWithError(event Event, err error) {
	if err == nil {
		helper.FinishWithCond(event, true)
	}

	event.AddErr(err)
	helper.FinishWithCond(event, false)
}
