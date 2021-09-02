// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEventHelper(t *testing.T) {
	helper := NewEventHelper(NewEventFactory())

	assert.NotNil(t, helper)
	assert.NotNil(t, helper.Factory)
}

func TestEventHelper_Start(t *testing.T) {
	helper := NewEventHelper(NewEventFactory())

	event := helper.Start("")
	assert.NotNil(t, event)
	assert.NotZero(t, event.GetStartTime().Unix())
}

func TestEventHelper_Finish(t *testing.T) {
	helper := NewEventHelper(NewEventFactory())

	event := helper.Start("")
	helper.Finish(event)

	assert.NotZero(t, event.GetStartTime().Unix())
	assert.NotZero(t, event.GetEndTime().Unix())
}

func TestEventHelper_FinishWithCond_WithSuccess(t *testing.T) {
	helper := NewEventHelper(NewEventFactory())

	event := helper.Start("")
	helper.FinishWithCond(event, true)

	assert.NotZero(t, event.GetStartTime().Unix())
	assert.NotZero(t, event.GetEndTime().Unix())
	assert.Equal(t, int64(1), event.GetCounter("success"))
}

func TestEventHelper_FinishWithCond_WithFailure(t *testing.T) {
	helper := NewEventHelper(NewEventFactory())

	event := helper.Start("")
	helper.FinishWithCond(event, false)

	assert.NotZero(t, event.GetStartTime().Unix())
	assert.NotZero(t, event.GetEndTime().Unix())
	assert.Equal(t, int64(1), event.GetCounter("failure"))
}

func TestEventHelper_FinishWithError_WithoutError(t *testing.T) {
	helper := NewEventHelper(NewEventFactory())

	event := helper.Start("")
	helper.FinishWithError(event, nil)

	assert.NotZero(t, event.GetStartTime().Unix())
	assert.NotZero(t, event.GetEndTime().Unix())
	assert.Equal(t, int64(1), event.GetCounter("success"))
	assert.Zero(t, event.GetErrCount(errors.New("")))
}

func TestEventHelper_FinishWithError_WithError(t *testing.T) {
	helper := NewEventHelper(NewEventFactory())

	event := helper.Start("")
	helper.FinishWithError(event, &MyErr{})

	assert.NotZero(t, event.GetStartTime().Unix())
	assert.NotZero(t, event.GetEndTime().Unix())
	assert.Equal(t, int64(1), event.GetCounter("failure"))
	assert.NotZero(t, event.GetErrCount(&MyErr{}))
}

type MyErr struct{}

func (err MyErr) Error() string {
	return ""
}
