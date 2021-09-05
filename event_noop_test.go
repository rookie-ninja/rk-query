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

func TestEventNoop_AllInOne(t *testing.T) {
	event := &eventNoop{}

	// No panic should be occur
	event.SetStartTime(time.Now())
	assert.NotNil(t, event.GetStartTime())
	event.SetEndTime(time.Now())
	assert.NotNil(t, event.GetEndTime())
	event.AddPayloads(zap.String("key", "value"))
	assert.NotNil(t, event.ListPayloads())
	assert.Empty(t, event.GetEventId())
	event.SetEventId("")
	assert.Empty(t, event.GetTraceId())
	event.SetTraceId("")
	assert.Empty(t, event.GetRequestId())
	event.SetRequestId("")
	event.AddErr(errors.New(""))
	assert.Zero(t, event.GetErrCount(errors.New("")))
	assert.Empty(t, event.GetOperation())
	event.SetOperation("")
	assert.Empty(t, event.GetRemoteAddr())
	event.SetRemoteAddr("")
	assert.Empty(t, event.GetResCode())
	event.SetResCode("")
	assert.Equal(t, NotStarted, event.GetEventStatus())
	event.StartTimer("")
	event.EndTimer("")
	event.UpdateTimerMs("", 0)
	event.UpdateTimerMsWithSample("", 0, 0)
	assert.Zero(t, event.GetTimeElapsedMs(""))
	assert.Empty(t, event.GetValueFromPair(""))
	event.AddPair("", "")
	assert.Zero(t, event.GetCounter(""))
	event.SetCounter("", 0)
	event.IncCounter("", 0)
	event.Finish()
}
