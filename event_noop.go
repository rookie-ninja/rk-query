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

// Set start timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventNoop) SetStartTime(time.Time) {
	// Noop
}

// Get start time of current event data.
func (event *eventNoop) GetStartTime() time.Time {
	return time.Now()
}

// Set end timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventNoop) SetEndTime(time.Time) {
	// Noop
}

// Get end time of current event data.
func (event *eventNoop) GetEndTime() time.Time {
	return time.Now()
}

// ************* Payload *************

// Add payload as zap.Field.
// Payload could be anything with RPC requests or user event such as http request param.
func (event *eventNoop) AddPayloads(...zap.Field) {
	// Noop
}

// List payloads.
func (event *eventNoop) ListPayloads() []zap.Field {
	return []zap.Field{}
}

// ************* Identity *************

// Get event id of current event.
func (event *eventNoop) GetEventId() string {
	return ""
}

// Set event id of current event.
// A new event id would be created while event data was created from EventFactory.
// User could override event id with this function.
func (event *eventNoop) SetEventId(string) {
	// Noop
}

// Get trace id of current event.
func (event *eventNoop) GetTraceId() string {
	return ""
}

// Set trace id of current event.
func (event *eventNoop) SetTraceId(string) {
	// Noop
}

// Get request id of current event.
func (event *eventNoop) GetRequestId() string {
	return ""
}

// Set request id of current event.
func (event *eventNoop) SetRequestId(string) {
	// Noop
}

// ************* Error *************

// Add an error into event which could be printed with error.Error() function.
func (event *eventNoop) AddErr(error) {
	// Noop
}

// Get error count.
// We will use value of error.Error() as the key.
func (event *eventNoop) GetErrCount(error) int64 {
	return 0
}

// ************* Event *************

// Get operation of current event.
func (event *eventNoop) GetOperation() string {
	return ""
}

// Set operation of current event.
func (event *eventNoop) SetOperation(string) {
	// Noop
}

// Get remote address of current event.
func (event *eventNoop) GetRemoteAddr() string {
	return ""
}

// Set remote address of current event, mainly used in RPC calls.
// Default value of <localhost> would be assigned while creating event via EventFactory.
func (event *eventNoop) SetRemoteAddr(string) {
	// Noop
}

// Get response code of current event.
// Mainly used in RPC calls.
func (event *eventNoop) GetResCode() string {
	return ""
}

// Set response code of current event.
func (event *eventNoop) SetResCode(string) {
	// Noop
}

// Get event status of current event.
// Available event status as bellow:
// 1: NotStarted
// 2: InProgress
// 3: Ended
func (event *eventNoop) GetEventStatus() eventStatus {
	return NotStarted
}

// Start timer of current sub event.
func (event *eventNoop) StartTimer(string) {
	// Noop
}

// End timer of current sub event.
func (event *eventNoop) EndTimer(string) {
	// Noop
}

// Update timer of current sub event with time elapsed in milli seconds.
func (event *eventNoop) UpdateTimerMs(string, int64) {
	// Noop
}

// Update timer of current sub event with time elapsed in milli seconds and sample.
func (event *eventNoop) UpdateTimerMsWithSample(string, int64, int64) {
	// Noop
}

// Get timer elapsed in milli seconds.
func (event *eventNoop) GetTimeElapsedMs(string) int64 {
	return 0
}

// Get value with key in pairs.
func (event *eventNoop) GetValueFromPair(string) string {
	return ""
}

// Add value with key in pairs.
func (event *eventNoop) AddPair(string, string) {
	// Noop
}

// Get counter of current event.
func (event *eventNoop) GetCounter(string) int64 {
	return 0
}

// Set counter of current event.
func (event *eventNoop) SetCounter(string, int64) {
	// Noop
}

// Increase counter of current event.
func (event *eventNoop) IncCounter(string, int64) {
	// Noop
}

// Set event status and flush to logger.
func (event *eventNoop) Finish() {
	// Noop
}
