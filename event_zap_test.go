// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkquery

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestEventStatus_String(t *testing.T) {
	assert.Equal(t, "NotStarted", NotStarted.String())
	assert.Equal(t, "InProgress", InProgress.String())
	assert.Equal(t, "Ended", Ended.String())
}

func TestEncoding_String(t *testing.T) {
	assert.Equal(t, "console", CONSOLE.String())
	assert.Equal(t, "json", JSON.String())
}

func TestToEncoding_HappyCase(t *testing.T) {
	assert.Equal(t, JSON, ToEncoding("json"))
	assert.Equal(t, CONSOLE, ToEncoding("console"))
}

func TestToEncoding_WithInvalidString(t *testing.T) {
	assert.Equal(t, CONSOLE, ToEncoding("invalid"))
}

func TestEventZap_SetStartTime(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	curr := time.Now()
	event.SetStartTime(curr)
	assert.Equal(t, curr, event.GetStartTime())
	assert.Equal(t, InProgress, event.GetEventStatus())
}

func TestEventZap_GetStartTime(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	curr := time.Now()
	event.SetStartTime(curr)
	assert.Equal(t, curr, event.GetStartTime())
}

func TestEventZap_SetEndTime(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetStartTime(time.Now())

	curr := time.Now()
	event.SetEndTime(curr)
	assert.Equal(t, curr, event.GetEndTime())
	assert.Equal(t, Ended, event.GetEventStatus())
}

func TestEventZap_GetEndTime(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetStartTime(time.Now())

	curr := time.Now()
	event.SetEndTime(curr)
	assert.Equal(t, curr, event.GetEndTime())
	assert.Equal(t, Ended, event.GetEventStatus())
}

func TestEventZap_AddPayloads(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.AddPayloads(zap.String("key", "value"))

	assert.Len(t, event.ListPayloads(), 1)
}

func TestEventZap_ListPayloads(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.AddPayloads(zap.String("key", "value"))

	assert.Len(t, event.ListPayloads(), 1)
}

func TestEventZap_SetEventId(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetEventId("ut-event")
	assert.Equal(t, "ut-event", event.GetEventId())
}

func TestEventZap_GetEventId(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetEventId("ut-event")
	assert.Equal(t, "ut-event", event.GetEventId())
}

func TestEventZap_SetTraceId(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetTraceId("ut-trace")
	assert.Equal(t, "ut-trace", event.GetTraceId())
}

func TestEventZap_GetTraceId(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetTraceId("ut-trace")
	assert.Equal(t, "ut-trace", event.GetTraceId())
}

func TestEventZap_SetRequestId(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetRequestId("ut-request")
	assert.Equal(t, "ut-request", event.GetRequestId())
}

func TestEventZap_GetRequestId(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetRequestId("ut-request")
	assert.Equal(t, "ut-request", event.GetRequestId())
}

func TestEventZap_AddErr_WithNilErr(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.AddErr(nil)
	assert.Empty(t, event.(*eventZap).errors.Fields)
}

func TestEventZap_AddErr_WithEmptyStrErr(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	err := errors.New("")
	event.AddErr(err)
	assert.Equal(t, int64(1), event.GetErrCount(err))
}

func TestEventZap_AddErr_HappyCase(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	err := errors.New("ut-err")
	event.AddErr(err)
	assert.Equal(t, int64(1), event.GetErrCount(err))

	event.AddErr(err)
	assert.Equal(t, int64(2), event.GetErrCount(err))
}

func TestEventZap_GetErrCount_WithNonExistErr(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.AddErr(errors.New("ut-err"))
	assert.Zero(t, event.GetErrCount(errors.New("non-exist-err")))
}

func TestEventZap_SetOperation(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetOperation("ut-operation")
	assert.Equal(t, "ut-operation", event.GetOperation())
}

func TestEventZap_GetOperation(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetOperation("ut-operation")
	assert.Equal(t, "ut-operation", event.GetOperation())
}

func TestEventZap_SetRemoteAddr(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetRemoteAddr("ut-remote-addr")
	assert.Equal(t, "ut-remote-addr", event.GetRemoteAddr())
}

func TestEventZap_GetRemoteAddr(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetRemoteAddr("ut-remote-addr")
	assert.Equal(t, "ut-remote-addr", event.GetRemoteAddr())
}

func TestEventZap_SetResCode(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetResCode("ut-res-code")
	assert.Equal(t, "ut-res-code", event.GetResCode())
}

func TestEventZap_GetResCode(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetResCode("ut-res-code")
	assert.Equal(t, "ut-res-code", event.GetResCode())
}

func TestEventZap_GetEventStatus(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.(*eventZap).status = InProgress
	assert.Equal(t, InProgress, event.GetEventStatus())
}

func TestEventZap_StartTimer_WithInvalidEventStatus(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	// event status is not InProgress, so nothing will happen
	event.StartTimer("ut-timer")
	assert.Equal(t, int64(-1), event.GetTimeElapsedMs("ut-timer"))
}

func TestEventZap_StartTimer_WithEmptyName(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.StartTimer("")
	assert.Equal(t, int64(-1), event.GetTimeElapsedMs(""))
}

func TestEventZap_StartTimer_WithCallMultipleTimes(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetStartTime(time.Now())
	event.StartTimer("ut-timer")
	event.StartTimer("ut-timer")

	assert.True(t, event.GetTimeElapsedMs("ut-timer") >= 0)
}

func TestEventZap_EndTimer_WithInvalidEventStatus(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	// event status is not InProgress, so nothing will happen
	event.EndTimer("ut-timer")
	assert.Equal(t, int64(-1), event.GetTimeElapsedMs("ut-timer"))
}

func TestEventZap_EndTimer_WithEmptyName(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.EndTimer("")
	assert.Equal(t, int64(-1), event.GetTimeElapsedMs(""))
}

func TestEventZap_EndTimer_WithoutStart(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.EndTimer("ut-timer")
	assert.Equal(t, int64(-1), event.GetTimeElapsedMs("ut-timer"))
}

func TestEventZap_EndTimer_HappyCase(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetStartTime(time.Now())
	event.StartTimer("ut-timer")
	event.EndTimer("ut-timer")
	assert.True(t, event.GetTimeElapsedMs("ut-timer") >= 0)
}

func TestEventZap_UpdateTimerMs(t *testing.T) {

}

func TestEventZap_UpdateTimerMs_WithInvalidEventStatus(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	// event status is not InProgress, so nothing will happen
	event.UpdateTimerMs("ut-timer", 1)
	assert.Equal(t, int64(-1), event.GetTimeElapsedMs("ut-timer"))
}

func TestEventZap_UpdateTimerMs_WithEmptyName(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.UpdateTimerMs("", 1)
	assert.Equal(t, int64(-1), event.GetTimeElapsedMs(""))
}

func TestEventZap_UpdateTimerMs_WithCallMultipleTimes(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetStartTime(time.Now())
	event.UpdateTimerMs("ut-timer", 1)
	event.UpdateTimerMs("ut-timer", 1)

	assert.True(t, event.GetTimeElapsedMs("ut-timer") >= 0)
}

func TestEventZap_UpdateTimerMsWithSample_WithInvalidEventStatus(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	// event status is not InProgress, so nothing will happen
	event.UpdateTimerMsWithSample("ut-timer", 1, 1)
	assert.Equal(t, int64(-1), event.GetTimeElapsedMs("ut-timer"))
}

func TestEventZap_UpdateTimerMsWithSample_WithEmptyName(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.UpdateTimerMsWithSample("", 1, 1)
	assert.Equal(t, int64(-1), event.GetTimeElapsedMs(""))
}

func TestEventZap_UpdateTimerMsWithSample_WithCallMultipleTimes(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetStartTime(time.Now())
	event.UpdateTimerMsWithSample("ut-timer", 1, 1)
	event.UpdateTimerMsWithSample("ut-timer", 1, 1)

	assert.True(t, event.GetTimeElapsedMs("ut-timer") >= 0)
}

func TestEventZap_AddPair(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.AddPair("key", "value")
	assert.Equal(t, "value", event.GetValueFromPair("key"))
}

func TestEventZap_SetCounter(t *testing.T) {
	event := NewEventFactory().CreateEvent()
	event.SetCounter("ut-counter", 1)
	assert.Equal(t, int64(1), event.GetCounter("ut-counter"))

	event.IncCounter("ut-counter", 1)
	assert.Equal(t, int64(2), event.GetCounter("ut-counter"))
}

func TestEventZap_Finish_WithQuiteMode(t *testing.T) {
	event := NewEventFactory().CreateEvent(WithQuietMode(true))
	// No output expected
	event.Finish()
}

func TestEventZap_Finish_WithJsonEncoding(t *testing.T) {
	event := NewEventFactory().CreateEvent(WithEncoding(JSON))
	// JSON format expected
	event.Finish()
}

func TestEventZap_Finish_WithConsoleEncoding(t *testing.T) {
	event := NewEventFactory().CreateEvent(WithEncoding(CONSOLE))
	// Console format expected
	event.Finish()
}

func TestEventZap_Finish_WithTracker(t *testing.T) {
	event := NewEventFactory().CreateEvent(WithEncoding(CONSOLE))
	event.SetStartTime(time.Now())
	event.StartTimer("ut-timer")
	event.EndTimer("ut-timer")
	event.SetEndTime(time.Now())
	// Console format expected
	event.Finish()
}
