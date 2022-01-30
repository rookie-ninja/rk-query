// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkquery

import (
	"go.uber.org/zap"
	"time"
)

type eventNoop struct{}

// ************* Time *************

// SetStartTime sets start timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventNoop) SetStartTime(time.Time) {
	// Noop
}

// GetStartTime returns start time of current event data.
func (event *eventNoop) GetStartTime() time.Time {
	return time.Now()
}

// SetEndTime sets end timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventNoop) SetEndTime(time.Time) {
	// Noop
}

// GetEndTime returns end time of current event data.
func (event *eventNoop) GetEndTime() time.Time {
	return time.Now()
}

// ************* Payload *************

// AddPayloads function add payload as zap.Field.
// Payload could be anything with RPC requests or user event such as http request param.
func (event *eventNoop) AddPayloads(...zap.Field) {
	// Noop
}

// ListPayloads will lists payloads.
func (event *eventNoop) ListPayloads() []zap.Field {
	return []zap.Field{}
}

// ************* Identity *************

// GetEventId returns event id of current event.
func (event *eventNoop) GetEventId() string {
	return ""
}

// SetEventId sets event id of current event.
// A new event id would be created while event data was created from EventFactory.
// User could override event id with this function.
func (event *eventNoop) SetEventId(string) {
	// Noop
}

// GetTraceId returns trace id of current event.
func (event *eventNoop) GetTraceId() string {
	return ""
}

// SetTraceId set trace id of current event.
func (event *eventNoop) SetTraceId(string) {
	// Noop
}

// GetRequestId returns request id of current event.
func (event *eventNoop) GetRequestId() string {
	return ""
}

// SetRequestId set request id of current event.
func (event *eventNoop) SetRequestId(string) {
	// Noop
}

// ************* Error *************

// AddErr function adds an error into event which could be printed with error.Error() function.
func (event *eventNoop) AddErr(error) {
	// Noop
}

// GetErrCount returns error count.
// We will use value of error.Error() as the key.
func (event *eventNoop) GetErrCount(error) int64 {
	return 0
}

// ************* Event *************

// GetOperation returns operation of current event.
func (event *eventNoop) GetOperation() string {
	return ""
}

// SetOperation sets operation of current event.
func (event *eventNoop) SetOperation(string) {
	// Noop
}

// GetRemoteAddr returns remote address of current event.
func (event *eventNoop) GetRemoteAddr() string {
	return ""
}

// SetRemoteAddr sets remote address of current event, mainly used in RPC calls.
// Default value of <localhost> would be assigned while creating event via EventFactory.
func (event *eventNoop) SetRemoteAddr(string) {
	// Noop
}

// GetResCode returns response code of current event.
// Mainly used in RPC calls.
func (event *eventNoop) GetResCode() string {
	return ""
}

// SetResCode sets response code of current event.
func (event *eventNoop) SetResCode(string) {
	// Noop
}

// GetEventStatus returns event status of current event.
// Available event status as bellow:
// 1: NotStarted
// 2: InProgress
// 3: Ended
func (event *eventNoop) GetEventStatus() eventStatus {
	return NotStarted
}

// StartTimer starts timer of current sub event.
func (event *eventNoop) StartTimer(string) {
	// Noop
}

// EndTimer ends timer of current sub event.
func (event *eventNoop) EndTimer(string) {
	// Noop
}

// UpdateTimerMs updates timer of current sub event with time elapsed in milli seconds.
func (event *eventNoop) UpdateTimerMs(string, int64) {
	// Noop
}

// UpdateTimerMsWithSample updates timer of current sub event with time elapsed in milli seconds and sample.
func (event *eventNoop) UpdateTimerMsWithSample(string, int64, int64) {
	// Noop
}

// GetTimeElapsedMs returns timer elapsed in milli seconds.
func (event *eventNoop) GetTimeElapsedMs(string) int64 {
	return 0
}

// GetValueFromPair returns value with key in pairs.
func (event *eventNoop) GetValueFromPair(string) string {
	return ""
}

// AddPair adds value with key in pairs.
func (event *eventNoop) AddPair(string, string) {
	// Noop
}

// GetCounter returns counter of current event.
func (event *eventNoop) GetCounter(string) int64 {
	return 0
}

// SetCounter sets counter of current event.
func (event *eventNoop) SetCounter(string, int64) {
	// Noop
}

// IncCounter increases counter of current event.
func (event *eventNoop) IncCounter(string, int64) {
	// Noop
}

// Finish sets event status and flush to logger.
func (event *eventNoop) Finish() {
	// Noop
}

// Sync flushes logs in buffer, mainly used for external syncer
func (event *eventNoop) Sync() {
	// Noop
}
