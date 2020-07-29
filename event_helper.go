// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"time"
)

// A helper function for easy use of EventData
type EventZapHelper struct {
	Factory *EventZapFactory
}

func NewEventZapHelper(factory *EventZapFactory) *EventZapHelper {
	if factory == nil {
		factory = NewEventZapFactory()
	}
	return &EventZapHelper{factory}
}

func (helper *EventZapHelper) Start(operation string) *EventZap {
	event := helper.Factory.CreateEventZap()

	event.SetOperation(operation)
	event.SetStartTime(time.Now())
	return event
}

func (helper *EventZapHelper) Finish(event *EventZap) {
	event.SetEndTime(time.Now())
	event.WriteLog()
}

func (helper *EventZapHelper) FinishWithCond(event *EventZap, success bool) {
	if success {
		event.SetCounter("success", 1)
	} else {
		event.SetCounter("failure", 1)
	}

	helper.Finish(event)
}

func (helper *EventZapHelper) FinishWithError(event *EventZap, err error) {
	if err == nil {
		helper.FinishWithCond(event, true)
	}

	event.AddErr(err)
	helper.FinishWithCond(event, false)
}
