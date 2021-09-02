// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"go.uber.org/zap"
	"time"
)

type Event interface {
	// ************* Time *************

	// Set start timer of current event. This can be overridden by user.
	// We keep this function open in order to mock event during unit test.
	SetStartTime(time.Time)

	// Get start time of current event data.
	GetStartTime() time.Time

	// Set end timer of current event. This can be overridden by user.
	// We keep this function open in order to mock event during unit test.
	SetEndTime(time.Time)

	// Get end time of current event data.
	GetEndTime() time.Time

	// ************* Payload *************

	// Add payload as zap.Field.
	// Payload could be anything with RPC requests or user event such as http request param.
	AddPayloads(...zap.Field)

	// List payloads.
	ListPayloads() []zap.Field

	// ************* Identity *************

	// Get event id of current event.
	GetEventId() string

	// Set event id of current event.
	// A new event id would be created while event data was created from EventFactory.
	// User could override event id with this function.
	SetEventId(string)

	// Get trace id of current event.
	GetTraceId() string

	// Set trace id of current event.
	SetTraceId(string)

	// Get request id of current event.
	GetRequestId() string

	// Set request id of current event.
	SetRequestId(string)

	// ************* Error *************

	// Add an error into event which could be printed with error.Error() function.
	AddErr(error)

	// Get error count.
	// We will use value of error.Error() as the key.
	GetErrCount(error) int64

	// ************* Event *************

	// Get operation of current event.
	GetOperation() string

	// Set operation of current event.
	SetOperation(string)

	// Get remote address of current event.
	GetRemoteAddr() string

	// Set remote address of current event, mainly used in RPC calls.
	// Default value of <localhost> would be assigned while creating event via EventFactory.
	SetRemoteAddr(string)

	// Get response code of current event.
	// Mainly used in RPC calls.
	GetResCode() string

	// Set response code of current event.
	SetResCode(string)

	// Get event status of current event.
	// Available event status as bellow:
	// 1: NotStarted
	// 2: InProgress
	// 3: Ended
	GetEventStatus() eventStatus

	// Start timer of current sub event.
	StartTimer(string)

	// End timer of current sub event.
	EndTimer(string)

	// Update timer of current sub event with time elapsed in milli seconds.
	UpdateTimerMs(string, int64)

	// Update timer of current sub event with time elapsed in milli seconds and sample.
	UpdateTimerMsWithSample(string, int64, int64)

	// Get timer elapsed in milli seconds.
	GetTimeElapsedMs(string) int64

	// Get value with key in pairs.
	GetValueFromPair(string) string

	// Add value with key in pairs.
	AddPair(string, string)

	// Get counter of current event.
	GetCounter(string) int64

	// Set counter of current event.
	SetCounter(string, int64)

	// Increase counter of current event.
	IncCounter(string, int64)

	// Set event status and flush to logger.
	Finish()
}
