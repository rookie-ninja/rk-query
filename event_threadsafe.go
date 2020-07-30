package rk_query

import (
	"go.uber.org/zap"
	"sync"
	"time"
)

type EventThreadSafe struct{
	delegate Event
	lock     *sync.Mutex
}

func (event *EventThreadSafe) GetValue(key string) string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetValue(key)
}

func (event *EventThreadSafe) GetAppName() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetAppName()
}

func (event *EventThreadSafe) GetEventId() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventId()
}

func (event *EventThreadSafe) SetEventId(id string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetEventId(id)
}

func (event *EventThreadSafe) GetHostname() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetHostname()
}

func (event *EventThreadSafe) GetLogger() *zap.Logger {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetLogger()
}

func (event *EventThreadSafe) GetOperation() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetOperation()
}

func (event *EventThreadSafe) SetOperation(operation string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetOperation(operation)
}

func (event *EventThreadSafe) GetEventStatus() eventStatus {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventStatus()
}

func (event *EventThreadSafe) SetStartTime(curr time.Time) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetStartTime(curr)
}

func (event *EventThreadSafe) GetStartTime() time.Time {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetStartTime()
}

func (event *EventThreadSafe) GetEndTime() time.Time {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEndTime()
}

func (event *EventThreadSafe) SetEndTime(curr time.Time) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetEndTime(curr)
}

func (event *EventThreadSafe) StartTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.StartTimer(name)
}

func (event *EventThreadSafe) EndTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.EndTimer(name)
}

func (event *EventThreadSafe) UpdateTimer(name string, elapsedMS int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.UpdateTimer(name, elapsedMS)
}

func (event *EventThreadSafe) UpdateTimerWithSample(name string, elapsedMS, sample int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.UpdateTimerWithSample(name, elapsedMS, sample)
}

func (event *EventThreadSafe) GetTimeElapsedMS(name string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.GetTimeElapsedMS(name)
}

func (event *EventThreadSafe) GetRemoteAddr() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetRemoteAddr()
}

func (event *EventThreadSafe) SetRemoteAddr(addr string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetRemoteAddr(addr)
}

func (event *EventThreadSafe) GetCounter(key string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetCounter(key)
}

func (event *EventThreadSafe) SetCounter(key string, value int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetCounter(key, value)
}

func (event *EventThreadSafe) InCCounter(key string, value int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.InCCounter(key, value)
}

func (event *EventThreadSafe) AddPair(key, value string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddPair(key, value)
}

func (event *EventThreadSafe) AddErr(err error) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddErr(err)
}

func (event *EventThreadSafe) GetErrCount(err error) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return  event.delegate.GetErrCount(err)
}

func (event *EventThreadSafe) AddFields(fields ...zap.Field) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddFields(fields...)
}

func (event *EventThreadSafe) GetFields() []zap.Field {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetFields()
}

func (event *EventThreadSafe) RecordHistoryEvent(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.RecordHistoryEvent(name)
}

func (event *EventThreadSafe) WriteLog() {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.WriteLog()
}

func (event *EventThreadSafe) setLogger(logger *zap.Logger) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setLogger(logger)
}

func (event *EventThreadSafe) setFormat(format Format) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setFormat(format)
}

func (event *EventThreadSafe) setQuietMode(quietMode bool) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setQuietMode(quietMode)
}

func (event *EventThreadSafe) setAppName(appName string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setAppName(appName)
}

func (event *EventThreadSafe) setHostname(hostname string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setHostname(hostname)
}
