package rk_query

import (
	"go.uber.org/zap"
	"sync"
	"time"
)

type eventThreadSafe struct{
	delegate Event
	lock     *sync.Mutex
}

func (event *eventThreadSafe) GetValue(key string) string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetValue(key)
}

func (event *eventThreadSafe) GetAppName() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetAppName()
}

func (event *eventThreadSafe) GetEventId() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventId()
}

func (event *eventThreadSafe) SetEventId(id string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetEventId(id)
}

func (event *eventThreadSafe) GetHostname() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetHostname()
}

func (event *eventThreadSafe) GetLogger() *zap.Logger {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetLogger()
}

func (event *eventThreadSafe) GetOperation() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetOperation()
}

func (event *eventThreadSafe) SetOperation(operation string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetOperation(operation)
}

func (event *eventThreadSafe) GetEventStatus() eventStatus {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEventStatus()
}

func (event *eventThreadSafe) SetStartTime(curr time.Time) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetStartTime(curr)
}

func (event *eventThreadSafe) GetStartTime() time.Time {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetStartTime()
}

func (event *eventThreadSafe) GetEndTime() time.Time {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetEndTime()
}

func (event *eventThreadSafe) SetEndTime(curr time.Time) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetEndTime(curr)
}

func (event *eventThreadSafe) StartTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.StartTimer(name)
}

func (event *eventThreadSafe) EndTimer(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.EndTimer(name)
}

func (event *eventThreadSafe) UpdateTimer(name string, elapsedMS int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.UpdateTimer(name, elapsedMS)
}

func (event *eventThreadSafe) UpdateTimerWithSample(name string, elapsedMS, sample int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.UpdateTimerWithSample(name, elapsedMS, sample)
}

func (event *eventThreadSafe) GetTimeElapsedMS(name string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.GetTimeElapsedMS(name)
}

func (event *eventThreadSafe) GetRemoteAddr() string {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetRemoteAddr()
}

func (event *eventThreadSafe) SetRemoteAddr(addr string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetRemoteAddr(addr)
}

func (event *eventThreadSafe) GetCounter(key string) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetCounter(key)
}

func (event *eventThreadSafe) SetCounter(key string, value int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.SetCounter(key, value)
}

func (event *eventThreadSafe) InCCounter(key string, value int64) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.InCCounter(key, value)
}

func (event *eventThreadSafe) AddPair(key, value string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddPair(key, value)
}

func (event *eventThreadSafe) AddErr(err error) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddErr(err)
}

func (event *eventThreadSafe) GetErrCount(err error) int64 {
	event.lock.Lock()
	defer event.lock.Unlock()

	return  event.delegate.GetErrCount(err)
}

func (event *eventThreadSafe) AddFields(fields ...zap.Field) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.AddFields(fields...)
}

func (event *eventThreadSafe) GetFields() []zap.Field {
	event.lock.Lock()
	defer event.lock.Unlock()

	return event.delegate.GetFields()
}

func (event *eventThreadSafe) RecordHistoryEvent(name string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.RecordHistoryEvent(name)
}

func (event *eventThreadSafe) WriteLog() {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.delegate.WriteLog()
}

func (event *eventThreadSafe) setLogger(logger *zap.Logger) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setLogger(logger)
}

func (event *eventThreadSafe) setFormat(format format) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setFormat(format)
}

func (event *eventThreadSafe) setQuietMode(quietMode bool) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setQuietMode(quietMode)
}

func (event *eventThreadSafe) setAppName(appName string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setAppName(appName)
}

func (event *eventThreadSafe) setHostname(hostname string) {
	event.lock.Lock()
	defer event.lock.Unlock()

	event.setHostname(hostname)
}
