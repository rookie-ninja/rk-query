// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

// Package rkquery can be used for creating Event instance for logging.
package rkquery

import (
	"go.uber.org/zap"
	"time"
)

// Event is used to record any event related stuff like RPC
type Event interface {
	// ************* Time *************

	// SetStartTime sets start timer of current event. This can be overridden by user.
	// We keep this function open in order to mock event during unit test.
	SetStartTime(time.Time)

	// GetStartTime returns start time of current event data.
	GetStartTime() time.Time

	// SetEndTime sets end timer of current event. This can be overridden by user.
	// We keep this function open in order to mock event during unit test.
	SetEndTime(time.Time)

	// GetEndTime returns end time of current event data.
	GetEndTime() time.Time

	// ************* Payload *************

	// AddPayloads function add payload as zap.Field.
	// Payload could be anything with RPC requests or user event such as http request param.
	AddPayloads(...zap.Field)

	// ListPayloads will lists payloads.
	ListPayloads() []zap.Field

	// ************* Identity *************

	// GetEventId returns event id of current event.
	GetEventId() string

	// SetEventId sets event id of current event.
	// A new event id would be created while event data was created from EventFactory.
	// User could override event id with this function.
	SetEventId(string)

	// GetTraceId returns trace id of current event.
	GetTraceId() string

	// SetTraceId set trace id of current event.
	SetTraceId(string)

	// GetRequestId returns request id of current event.
	GetRequestId() string

	// SetRequestId set request id of current event.
	SetRequestId(string)

	// ************* Error *************

	// AddErr function adds an error into event which could be printed with error.Error() function.
	AddErr(error)

	// GetErrCount returns error count.
	// We will use value of error.Error() as the key.
	GetErrCount(error) int64

	// ************* Event *************

	// GetOperation returns operation of current event.
	GetOperation() string

	// SetOperation sets operation of current event.
	SetOperation(string)

	// GetRemoteAddr returns remote address of current event.
	GetRemoteAddr() string

	// SetRemoteAddr sets remote address of current event, mainly used in RPC calls.
	// Default value of <localhost> would be assigned while creating event via EventFactory.
	SetRemoteAddr(string)

	// GetResCode returns response code of current event.
	// Mainly used in RPC calls.
	GetResCode() string

	// SetResCode sets response code of current event.
	SetResCode(string)

	// GetEventStatus returns event status of current event.
	// Available event status as bellow:
	// 1: NotStarted
	// 2: InProgress
	// 3: Ended
	GetEventStatus() eventStatus

	// StartTimer starts timer of current sub event.
	StartTimer(string)

	// EndTimer ends timer of current sub event.
	EndTimer(string)

	// UpdateTimerMs updates timer of current sub event with time elapsed in milli seconds.
	UpdateTimerMs(string, int64)

	// UpdateTimerMsWithSample updates timer of current sub event with time elapsed in milli seconds and sample.
	UpdateTimerMsWithSample(string, int64, int64)

	// GetTimeElapsedMs returns timer elapsed in milli seconds.
	GetTimeElapsedMs(string) int64

	// GetValueFromPair returns value with key in pairs.
	GetValueFromPair(string) string

	// AddPair adds value with key in pairs.
	AddPair(string, string)

	// GetCounter returns counter of current event.
	GetCounter(string) int64

	// SetCounter sets counter of current event.
	SetCounter(string, int64)

	// IncCounter increases counter of current event.
	IncCounter(string, int64)

	// Finish sets event status and flush to logger.
	Finish()

	// Sync flushes logger in buffer
	Sync()
}
