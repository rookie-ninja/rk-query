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

func TestEventThreadSafe_AllInOne(t *testing.T) {
	event := NewEventFactory().CreateEventThreadSafe(WithQuietMode(true))

	curr := time.Now()

	// No panics should occur
	// Start time
	event.SetStartTime(curr)
	assert.Equal(t, curr, event.GetStartTime())
	// Payloads
	event.AddPayloads(zap.String("key", "value"))
	assert.Len(t, event.ListPayloads(), 1)
	// EventId
	event.SetEventId("ut-eventId")
	assert.Equal(t, "ut-eventId", event.GetEventId())
	// TraceId
	event.SetTraceId("ut-trace")
	assert.Equal(t, "ut-trace", event.GetTraceId())
	// RequestId
	event.SetRequestId("ut-request")
	assert.Equal(t, "ut-request", event.GetRequestId())
	// Error
	utErr := errors.New("")
	event.AddErr(utErr)
	assert.Equal(t, int64(1), event.GetErrCount(utErr))
	// Operation
	event.SetOperation("ut-operation")
	assert.Equal(t, "ut-operation", event.GetOperation())
	// Remote address
	event.SetRemoteAddr("ut-remote-addr")
	assert.Equal(t, "ut-remote-addr", event.GetRemoteAddr())
	// ResCode
	event.SetResCode("ut-res-code")
	assert.Equal(t, "ut-res-code", event.GetResCode())
	// Event status
	assert.True(t, event.GetEventStatus() > 0)
	// Timer
	event.StartTimer("ut-timer")
	event.UpdateTimerMs("ut-timer", 1)
	event.UpdateTimerMsWithSample("ut-timer", 1, 1)
	event.EndTimer("ut-timer")
	assert.True(t, event.GetTimeElapsedMs("ut-timer") >= int64(0))
	// Pair
	event.AddPair("key", "value")
	assert.Equal(t, "value", event.GetValueFromPair("key"))
	// Counter
	event.SetCounter("ut-counter", 1)
	event.IncCounter("ut-counter", 1)
	assert.Equal(t, int64(2), event.GetCounter("ut-counter"))
	// EndTime
	event.SetEndTime(curr)
	assert.Equal(t, curr, event.GetEndTime())
	event.Finish()
}
