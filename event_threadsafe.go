// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"go.uber.org/zap"
	"sync"
)

type ThreadSafeEventImpl struct {
	delegate Event
	lock     *sync.Mutex
}

func NewThreadSafeEventImpl(delegate Event) *ThreadSafeEventImpl {
	event := ThreadSafeEventImpl{}
	event.delegate = delegate
	event.lock = &sync.Mutex{}
	return &event
}

func (event *ThreadSafeEventImpl) GetAppName() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetAppName()
}

func (event *ThreadSafeEventImpl) GetHostName() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.GetHostName()
}

func (event *ThreadSafeEventImpl) GetZapLogger() *zap.Logger {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetZapLogger()
}

func (event *ThreadSafeEventImpl) GetOperation() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetOperation()
}

func (event *ThreadSafeEventImpl) SetOperation(operation string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetOperation(operation)
}

func (event *ThreadSafeEventImpl) Reset() {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.Reset()
}

func (event *ThreadSafeEventImpl) GetEventStatus() eventDataStatus {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventStatus()
}

func (event *ThreadSafeEventImpl) GetEndTimeMS() int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEndTimeMS()
}

func (event *ThreadSafeEventImpl) GetStartTimeMS() int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetStartTimeMS()
}

func (event *ThreadSafeEventImpl) SetStartTimeMS(nowMS int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetStartTimeMS(nowMS)
}

func (event *ThreadSafeEventImpl) SetEndTimeMS(endTime int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetEndTimeMS(endTime)
}

func (event *ThreadSafeEventImpl) StartTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.StartTimer(name)
}

func (event *ThreadSafeEventImpl) EndTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.EndTimer(name)
}

func (event *ThreadSafeEventImpl) UpdateTimer(name string, elapsedMS int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.UpdateTimer(name, elapsedMS)
}

func (event *ThreadSafeEventImpl) UpdateTimerWithSample(name string, elapsedMS, sample int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.UpdateTimerWithSample(name, elapsedMS, sample)
}

func (event *ThreadSafeEventImpl) GetTimeElapsedMS(timerName string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetTimeElapsedMS(timerName)
}

func (event *ThreadSafeEventImpl) GetRemoteAddr() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetRemoteAddr()
}

func (event *ThreadSafeEventImpl) SetRemoteAddr(addr string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetRemoteAddr(addr)
}

func (event *ThreadSafeEventImpl) GetCounter(name string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetCounter(name)
}

func (event *ThreadSafeEventImpl) SetCounter(name string, value int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetCounter(name, value)
}

func (event *ThreadSafeEventImpl) InCCounter(name string, delta int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.InCCounter(name, delta)
}

func (event *ThreadSafeEventImpl) AddKv(key, value string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddKv(key, value)
}

func (event *ThreadSafeEventImpl) AddErr(err error) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddErr(err)
}

func (event *ThreadSafeEventImpl) AppendKv(key, value string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AppendKv(key, value)
}

func (event *ThreadSafeEventImpl) GetValue(key string) string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.GetValue(key)
}

func (event *ThreadSafeEventImpl) FinishCurrentEvent(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.FinishCurrentEvent(name)
}

func (event *ThreadSafeEventImpl) RecordHistoryEvent(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.RecordHistoryEvent(name)
}

func (event *ThreadSafeEventImpl) WriteLog() {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.WriteLog()
}

func (event *ThreadSafeEventImpl) ToZapFieldsMin() []zap.Field {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.ToZapFieldsMin()
}

func (event *ThreadSafeEventImpl) ToZapFields() []zap.Field {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.ToZapFields()
}

func (event *ThreadSafeEventImpl) GetEventHistory() *eventHistory {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventHistory()
}