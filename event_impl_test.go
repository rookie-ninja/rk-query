// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNewEventImpl(t *testing.T) {
	event := NewEventImpl(
		&RealTimeSource{},
		"app",
		"host",
		make(map[string]string),
		zap.NewNop(),
		RK,
		false,
		make([]eventEntryListener, 0),
		false)

	assert.NotNil(t, event)
	assert.NotNil(t, event.timeSource)
	assert.NotEmpty(t, event.appName)
	assert.NotEmpty(t, event.hostName)
	assert.NotNil(t, event.zapLogger)
	assert.Empty(t, event.listeners)
	assert.Empty(t, event.defaultKvs)
	assert.False(t, event.quietMode)
	assert.Equal(t, RK, event.format)
	assert.Empty(t, event.tracker)
	assert.Empty(t, event.kvs)
	assert.Empty(t, event.counters)
	assert.Empty(t, event.errors)
	assert.Equal(t, notStarted, event.status)
	assert.NotNil(t, event.mutex)
}

func TestEventImpl_GetAppName(t *testing.T) {
	assert.Equal(t, "app", getEventImpl().GetAppName())
}

func TestEventImpl_GetHostName(t *testing.T) {
	assert.Equal(t, "host", getEventImpl().GetHostName())
}

func TestEventImpl_GetZapLogger(t *testing.T) {
	assert.NotNil(t, getEventImpl().GetZapLogger())
}

func TestEventImpl_GetOperation(t *testing.T) {
	assert.Equal(t, Unknown, getEventImpl().GetOperation())
}

func TestEventImpl_SetOperation(t *testing.T) {
	event := getEventImpl()

	event.SetOperation("operation")
	assert.Equal(t, "operation", event.GetOperation())
}

func TestEventImpl_Reset(t *testing.T) {
	event := getEventImpl()

	assert.Equal(t, notStarted, event.status)
	assert.Zero(t, event.startTimeMS)
	assert.Zero(t, event.endTimeMS)
	assert.Empty(t, event.counters)
	assert.Empty(t, event.kvs)
	assert.Empty(t, event.tracker)
}

func TestEventImpl_GetEventStatus(t *testing.T) {
	assert.Equal(t, notStarted, getEventImpl().GetEventStatus())
}

func TestEventImpl_SetStartTimeMS(t *testing.T) {
	event := getEventImpl()

	event.SetStartTimeMS(1)
	assert.Equal(t, int64(1), event.startTimeMS)
	assert.Equal(t, inProgress, event.status)
}

func TestEventImpl_GetStartTimeMS(t *testing.T) {
	event := getEventImpl()

	event.SetStartTimeMS(1)
	assert.Equal(t, int64(1), event.GetStartTimeMS())
}

func TestEventImpl_SetEndTimeMS_WithoutStart(t *testing.T) {
	event := getEventImpl()

	event.SetEndTimeMS(1)
	assert.Zero(t, event.endTimeMS)
	assert.Equal(t, notStarted, event.status)
}

func TestEventImpl_GetEndTimeMS(t *testing.T) {
	event := getEventImpl()

	event.SetStartTimeMS(1)
	event.SetEndTimeMS(2)
	assert.Equal(t, ended, event.status)
	assert.Equal(t, int64(2), event.GetEndTimeMS())
}

func TestEventImpl_StartTimer_WithEmptyName(t *testing.T) {
	event := getEventImpl()
	event.SetStartTimeMS(1)
	event.StartTimer("")

	assert.Empty(t, event.tracker)
}

func TestEventImpl_StartTimer_WithoutStart(t *testing.T) {
	event := getEventImpl()
	event.StartTimer("")

	assert.Empty(t, event.tracker)
}

func TestEventImpl_StartTimer(t *testing.T) {
	event := getEventImpl()
	event.SetStartTimeMS(1)
	event.StartTimer("event")

	assert.NotEmpty(t, event.tracker)
}

func TestEventImpl_EndTimer_WithEmptyName(t *testing.T) {
	event := getEventImpl()
	event.SetStartTimeMS(1)
	event.EndTimer("")

	assert.Empty(t, event.tracker)
}

func TestEventImpl_EndTimer_WithoutStart(t *testing.T) {
	event := getEventImpl()
	event.EndTimer("")

	assert.Empty(t, event.tracker)
}

func TestEventImpl_EndTimer(t *testing.T) {
	event := getEventImpl()
	event.SetStartTimeMS(1)
	event.StartTimer("event")
	event.EndTimer("event")

	assert.NotEmpty(t, event.tracker)
}

func TestEventImpl_UpdateTimer_WithEmptyName(t *testing.T) {
	event := getEventImpl()
	event.SetStartTimeMS(1)
	event.UpdateTimer("", 1)

	assert.Empty(t, event.tracker)
}

func TestEventImpl_UpdateTimer_WithoutStart(t *testing.T) {
	event := getEventImpl()
	event.UpdateTimer("", 1)

	assert.Empty(t, event.tracker)
}

func TestEventImpl_UpdateTimer(t *testing.T) {
	event := getEventImpl()
	event.SetStartTimeMS(1)
	event.UpdateTimer("event", 1)

	assert.NotEmpty(t, event.tracker)
}

func TestEventImpl_GetTimeElapsedMS_WithoutTimer(t *testing.T) {
	event := getEventImpl()

	assert.Empty(t, event.tracker)
	assert.Equal(t, int64(-1), event.GetTimeElapsedMS(""))
}

func TestEventImpl_GetTimeElapsedMS(t *testing.T) {
	event := getEventImpl()
	event.SetStartTimeMS(1)
	event.StartTimer("event")
	event.EndTimer("event")

	assert.NotEmpty(t, event.tracker)
	assert.True(t, event.GetTimeElapsedMS("event") >= 0)
}

func TestEventImpl_GetRemoteAddr(t *testing.T) {
	assert.Empty(t, getEventImpl().GetRemoteAddr())
}

func TestEventImpl_SetRemoteAddr(t *testing.T) {
	event := getEventImpl()
	event.SetRemoteAddr("addr")
	assert.Equal(t, "addr", event.GetRemoteAddr())
}

func TestEventImpl_GetCounter_WithEmptyCounter(t *testing.T) {
	assert.Equal(t, int64(-1), getEventImpl().GetCounter("counter"))
}

func TestEventImpl_GetCounter(t *testing.T) {
	event := getEventImpl()
	event.SetCounter("counter", 1)
	assert.Equal(t, int64(1), event.GetCounter("counter"))
}

func TestEventImpl_SetCounter_WithEmptyName(t *testing.T) {
	event := getEventImpl()
	event.SetCounter("", 1)
	assert.Empty(t, event.counters)
}

func TestEventImpl_SetCounter(t *testing.T) {
	event := getEventImpl()
	event.SetCounter("counter", 1)
	assert.Equal(t, int64(1), event.GetCounter("counter"))
}

func TestEventImpl_InCCounter_WithEmptyName(t *testing.T) {
	event := getEventImpl()
	event.InCCounter("", 1)
	assert.Empty(t, event.counters)
}

func TestEventImpl_InCCounter_WithExistingCounter(t *testing.T) {
	event := getEventImpl()
	event.SetCounter("counter", 1)
	event.InCCounter("counter", 1)
	assert.Equal(t, int64(2), event.GetCounter("counter"))
}

func TestEventImpl_AddKv_WithEmptyName(t *testing.T) {
	event := getEventImpl()
	event.AddKv("", "value")
	assert.Empty(t, event.kvs)
}

func TestEventImpl_AddKv_WithEmptyValue(t *testing.T) {
	event := getEventImpl()
	event.AddKv("name", "")
	assert.Empty(t, event.kvs)
}

func TestEventImpl_AddKv(t *testing.T) {
	event := getEventImpl()
	event.AddKv("name", "value")
	assert.NotEmpty(t, event.kvs)
}

func TestEventImpl_AddErr(t *testing.T) {
	event := getEventImpl()
	event.AddErr(errors.New(""))
	assert.NotEmpty(t, event.errors)
}

func TestEventImpl_AddErr_WithMultiple(t *testing.T) {
	event := getEventImpl()
	event.AddErr(errors.New(""))
	event.AddErr(errors.New(""))

	assert.NotEmpty(t, event.errors)
	assert.Equal(t, 1, len(event.errors))
}

func TestEventImpl_GetErrCount(t *testing.T) {
	event := getEventImpl()
	event.AddErr(errors.New(""))
	event.AddErr(errors.New(""))

	assert.NotEmpty(t, event.errors)
	assert.Equal(t, int64(2), event.GetErrCount(errors.New("")))
}

func TestEventImpl_AppendKv_WithEmptyName(t *testing.T) {
	event := getEventImpl()
	event.AppendKv("", "value")
	assert.Empty(t, event.kvs)
}

func TestEventImpl_AppendKv_WithEmptyValue(t *testing.T) {
	event := getEventImpl()
	event.AppendKv("name", "")
	assert.Empty(t, event.kvs)
}

func TestEventImpl_AppendKv_WithExistingKey(t *testing.T) {
	event := getEventImpl()
	event.AddKv("name", "value")
	event.AppendKv("name", "value")
	assert.NotEmpty(t, event.kvs)
	assert.Equal(t, "value,value", event.GetValue("name"))
}

func TestEventImpl_AppendKv(t *testing.T) {
	event := getEventImpl()
	event.AddKv("name", "value")
	assert.NotEmpty(t, event.kvs)
	assert.Equal(t, "value", event.GetValue("name"))
}

func TestEventImpl_GetValue(t *testing.T) {
	event := getEventImpl()
	event.AddKv("name", "value")
	assert.NotEmpty(t, event.kvs)
	assert.Equal(t, "value", event.GetValue("name"))
}

func TestEventImpl_FinishCurrentTimer_WithoutStart(t *testing.T) {
	event := getEventImpl()
	event.FinishCurrentTimer("name")
	assert.Empty(t, event.tracker)
}

func TestEventImpl_FinishCurrentTimer(t *testing.T) {
	event := getEventImpl()
	event.SetStartTimeMS(1)
	event.StartTimer("name")
	event.FinishCurrentTimer("name")
	assert.NotEmpty(t, event.tracker)
}

func getEventImpl() *EventImpl {
	return NewEventImpl(
		&RealTimeSource{},
		"app",
		"host",
		make(map[string]string),
		zap.NewNop(),
		RK,
		false,
		make([]eventEntryListener, 0),
		false)
}