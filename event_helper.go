// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"time"
)

// A helper function for easy use of EventData
type EventHelper struct {
	Factory *EventFactory
}

func NewEventHelper(factory *EventFactory) *EventHelper {
	if factory == nil {
		factory = NewEventFactory()
	}
	return &EventHelper{factory}
}

func (helper *EventHelper) Start(operation string) Event {
	event := helper.Factory.CreateEvent()

	event.SetOperation(operation)
	event.SetStartTime(time.Now())
	return event
}

func (helper *EventHelper) Finish(event Event) {
	event.SetEndTime(time.Now())
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
