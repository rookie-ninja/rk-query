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

func TestNewEventHelperWithZapLogger(t *testing.T) {
	helper := NewEventHelperWithZapLogger("", &RealTimeSource{}, zap.NewNop())

	assert.NotNil(t, helper)
	assert.NotNil(t, helper.Factory)
	assert.NotNil(t, helper.TimeSource)
}

func TestEventHelper_Start(t *testing.T) {
	helper := NewEventHelperWithZapLogger("", &RealTimeSource{}, zap.NewNop())

	event := helper.Start("")
	assert.NotNil(t, event)
	assert.NotZero(t, event.GetStartTimeMS())
}

func TestEventHelper_Finish(t *testing.T) {
	helper := NewEventHelperWithZapLogger("", &RealTimeSource{}, zap.NewNop())

	event := helper.Start("")
	helper.Finish(event)

	assert.NotZero(t, event.GetStartTimeMS())
	assert.NotZero(t, event.GetEndTimeMS())
}

func TestEventHelper_FinishWithCond_WithSuccess(t *testing.T) {
	helper := NewEventHelperWithZapLogger("", &RealTimeSource{}, zap.NewNop())

	event := helper.Start("")
	helper.FinishWithCond(event, true)

	assert.NotZero(t, event.GetStartTimeMS())
	assert.NotZero(t, event.GetEndTimeMS())
	assert.Equal(t, int64(1), event.GetCounter("success"))
}

func TestEventHelper_FinishWithCond_WithFailure(t *testing.T) {
	helper := NewEventHelperWithZapLogger("", &RealTimeSource{}, zap.NewNop())

	event := helper.Start("")
	helper.FinishWithCond(event, false)

	assert.NotZero(t, event.GetStartTimeMS())
	assert.NotZero(t, event.GetEndTimeMS())
	assert.Equal(t, int64(1), event.GetCounter("failure"))
}

func TestEventHelper_FinishWithError_WithoutError(t *testing.T) {
	helper := NewEventHelperWithZapLogger("", &RealTimeSource{}, zap.NewNop())

	event := helper.Start("")
	helper.FinishWithError(event, nil)

	assert.NotZero(t, event.GetStartTimeMS())
	assert.NotZero(t, event.GetEndTimeMS())
	assert.Equal(t, 1, event.GetCounter("success"))
	assert.Zero(t, event.GetErrCount(errors.New("")))
}

func TestEventHelper_FinishWithError_WithError(t *testing.T) {
	helper := NewEventHelperWithZapLogger("", &RealTimeSource{}, zap.NewNop())

	event := helper.Start("")
	helper.FinishWithError(event, &MyErr{})

	assert.NotZero(t, event.GetStartTimeMS())
	assert.NotZero(t, event.GetEndTimeMS())
	assert.Equal(t, 1, event.GetCounter("failure"))
	assert.Zero(t, event.GetErrCount(&MyErr{}))
}
