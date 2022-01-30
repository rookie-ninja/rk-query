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

// SetStartTime sets start timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventThreadSafe) SetStartTime(curr time.Time) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetStartTime(curr)
}

// GetStartTime returns start time of current event data.
func (event *eventThreadSafe) GetStartTime() time.Time {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetStartTime()
}

// SetEndTime sets end timer of current event. This can be overridden by user.
// We keep this function open in order to mock event during unit test.
func (event *eventThreadSafe) SetEndTime(curr time.Time) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetEndTime(curr)
}

// GetEndTime returns end time of current event data.
func (event *eventThreadSafe) GetEndTime() time.Time {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEndTime()
}

// ************* Payload *************

// AddPayloads function add payload as zap.Field.
// Payload could be anything with RPC requests or user event such as http request param.
func (event *eventThreadSafe) AddPayloads(fields ...zap.Field) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddPayloads(fields...)
}

// ListPayloads will lists payloads.
func (event *eventThreadSafe) ListPayloads() []zap.Field {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.ListPayloads()
}

// ************* Identity *************

// GetEventId returns event id of current event.
func (event *eventThreadSafe) GetEventId() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventId()
}

// SetEventId sets event id of current event.
// A new event id would be created while event data was created from EventFactory.
// User could override event id with this function.
func (event *eventThreadSafe) SetEventId(id string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetEventId(id)
}

// GetTraceId returns trace id of current event.
func (event *eventThreadSafe) GetTraceId() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetTraceId()
}

// SetTraceId set trace id of current event.
func (event *eventThreadSafe) SetTraceId(id string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetTraceId(id)
}

// GetRequestId returns request id of current event.
func (event *eventThreadSafe) GetRequestId() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetRequestId()
}

// SetRequestId set request id of current event.
func (event *eventThreadSafe) SetRequestId(id string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetRequestId(id)
}

// ************* Error *************

// AddErr function adds an error into event which could be printed with error.Error() function.
func (event *eventThreadSafe) AddErr(err error) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddErr(err)
}

// GetErrCount returns error count.
// We will use value of error.Error() as the key.
func (event *eventThreadSafe) GetErrCount(err error) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetErrCount(err)
}

// ************* Event *************

// GetOperation returns operation of current event.
func (event *eventThreadSafe) GetOperation() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetOperation()
}

// SetOperation sets operation of current event.
func (event *eventThreadSafe) SetOperation(operation string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetOperation(operation)
}

// GetRemoteAddr returns remote address of current event.
func (event *eventThreadSafe) GetRemoteAddr() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetRemoteAddr()
}

// SetRemoteAddr sets remote address of current event, mainly used in RPC calls.
// Default value of <localhost> would be assigned while creating event via EventFactory.
func (event *eventThreadSafe) SetRemoteAddr(addr string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetRemoteAddr(addr)
}

// GetResCode returns response code of current event.
// Mainly used in RPC calls.
func (event *eventThreadSafe) GetResCode() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetResCode()
}

// SetResCode sets response code of current event.
func (event *eventThreadSafe) SetResCode(resCode string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetResCode(resCode)
}

// GetEventStatus returns event status of current event.
// Available event status as bellow:
// 1: NotStarted
// 2: InProgress
// 3: Ended
func (event *eventThreadSafe) GetEventStatus() eventStatus {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventStatus()
}

// StartTimer starts timer of current sub event.
func (event *eventThreadSafe) StartTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.StartTimer(name)
}

// EndTimer ends timer of current sub event.
func (event *eventThreadSafe) EndTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.EndTimer(name)
}

// UpdateTimerMs updates timer of current sub event with time elapsed in milli seconds.
func (event *eventThreadSafe) UpdateTimerMs(name string, elapsedMs int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.UpdateTimerMs(name, elapsedMs)
}

// UpdateTimerMsWithSample updates timer of current sub event with time elapsed in milli seconds.
func (event *eventThreadSafe) UpdateTimerMsWithSample(name string, elapsedMs, sample int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.UpdateTimerMsWithSample(name, elapsedMs, sample)
}

// GetTimeElapsedMs returns timer elapsed in milli seconds.
func (event *eventThreadSafe) GetTimeElapsedMs(name string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetTimeElapsedMs(name)
}

// GetValueFromPair returns value with key in pairs.
func (event *eventThreadSafe) GetValueFromPair(key string) string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetValueFromPair(key)
}

// AddPair adds value with key in pairs.
func (event *eventThreadSafe) AddPair(key, value string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddPair(key, value)
}

// GetCounter returns counter of current event.
func (event *eventThreadSafe) GetCounter(key string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetCounter(key)
}

// SetCounter sets counter of current event.
func (event *eventThreadSafe) SetCounter(key string, value int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetCounter(key, value)
}

// IncCounter increases counter of current event.
func (event *eventThreadSafe) IncCounter(key string, delta int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.IncCounter(key, delta)
}

// Finish sets event status and flush to logger.
func (event *eventThreadSafe) Finish() {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.Finish()
}

// Sync flushes logs in buffer, mainly used for external syncer
func (event *eventThreadSafe) Sync() {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.Sync()
}
