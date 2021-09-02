// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"go.uber.org/zap"
	"sync"
	"time"
)

type eventThreadSafe struct {
	delegate *eventZap
	lock     *sync.Mutex
}

// ************* Time *************

// Set start timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventThreadSafe) SetStartTime(curr time.Time) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetStartTime(curr)
}

// Get start time of current event data.
func (event *eventThreadSafe) GetStartTime() time.Time {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetStartTime()
}

// Set end timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventThreadSafe) SetEndTime(curr time.Time) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetEndTime(curr)
}

// Get end time of current event data.
func (event *eventThreadSafe) GetEndTime() time.Time {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEndTime()
}

// ************* Payload *************

// Add payload as zap.Field.
// Payload could be anything with RPC requests or user event such as http request param.
func (event *eventThreadSafe) AddPayloads(fields ...zap.Field) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddPayloads(fields...)
}

// List payloads.
func (event *eventThreadSafe) ListPayloads() []zap.Field {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.ListPayloads()
}

// ************* Identity *************

// Get event id of current event.
func (event *eventThreadSafe) GetEventId() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventId()
}

// Set event id of current event.
// A new event id would be created while event data was created from EventFactory.
// User could override event id with this function.
func (event *eventThreadSafe) SetEventId(id string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetEventId(id)
}

// Get trace id of current event.
func (event *eventThreadSafe) GetTraceId() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.GetTraceId()
}

// Set trace id of current event.
func (event *eventThreadSafe) SetTraceId(id string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetTraceId(id)
}

// Get request id of current event.
func (event *eventThreadSafe) GetRequestId() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetRequestId()
}

// Set request id of current event.
func (event *eventThreadSafe) SetRequestId(id string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetRequestId(id)
}

// ************* Error *************

// Add an error into event which could be printed with error.Error() function.
func (event *eventThreadSafe) AddErr(err error) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddErr(err)
}

// Get error count.
// We will use value of error.Error() as the key.
func (event *eventThreadSafe) GetErrCount(err error) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetErrCount(err)
}

// ************* Event *************

// Get operation of current event.
func (event *eventThreadSafe) GetOperation() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetOperation()
}

// Set operation of current event.
func (event *eventThreadSafe) SetOperation(operation string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetOperation(operation)
}

// Get remote address of current event.
func (event *eventThreadSafe) GetRemoteAddr() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetRemoteAddr()
}

// Set remote address of current event, mainly used in RPC calls.
// Default value of <localhost> would be assigned while creating event via EventFactory.
func (event *eventThreadSafe) SetRemoteAddr(addr string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetRemoteAddr(addr)
}

// Get response code of current event.
// Mainly used in RPC calls.
func (event *eventThreadSafe) GetResCode() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetResCode()
}

// Set response code of current event.
func (event *eventThreadSafe) SetResCode(resCode string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetResCode(resCode)
}

// Get event status of current event.
// Available event status as bellow:
// 1: NotStarted
// 2: InProgress
// 3: Ended
func (event *eventThreadSafe) GetEventStatus() eventStatus {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventStatus()
}

// Start timer of current sub event.
func (event *eventThreadSafe) StartTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.StartTimer(name)
}

// End timer of current sub event.
func (event *eventThreadSafe) EndTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.EndTimer(name)
}

// Update timer of current sub event with time elapsed in milli seconds.
func (event *eventThreadSafe) UpdateTimerMs(name string, elapsedMs int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.UpdateTimerMs(name, elapsedMs)
}

// Update timer of current sub event with time elapsed in milli seconds and sample.
func (event *eventThreadSafe) UpdateTimerMsWithSample(name string, elapsedMs, sample int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.UpdateTimerMsWithSample(name, elapsedMs, sample)
}

// Get timer elapsed in milli seconds.
func (event *eventThreadSafe) GetTimeElapsedMs(name string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetTimeElapsedMs(name)
}

// Get value with key in pairs.
func (event *eventThreadSafe) GetValueFromPair(key string) string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetValueFromPair(key)
}

// Add value with key in pairs.
func (event *eventThreadSafe) AddPair(key, value string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddPair(key, value)
}

// Get counter of current event.
func (event *eventThreadSafe) GetCounter(key string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetCounter(key)
}

// Set counter of current event.
func (event *eventThreadSafe) SetCounter(key string, value int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetCounter(key, value)
}

// Increase counter of current event.
func (event *eventThreadSafe) IncCounter(key string, delta int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.IncCounter(key, delta)
}

// Set event status and flush to logger.
func (event *eventThreadSafe) Finish() {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.Finish()
}
