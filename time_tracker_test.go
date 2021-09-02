// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTimeTrackerWithNilName(t *testing.T) {
	assert.Nil(t, newTimeTracker(""))
}

func TestNewTimeTrackerHappyCase(t *testing.T) {
	tracker := newTimeTracker("fake")
	assert.NotNil(t, tracker)
	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestGetNameHappyCase(t *testing.T) {
	assert.Equal(t, "fake", newTimeTracker("fake").name)
}

func TestStartWithNegativeNowMS(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(-1)
	// tracker should do nothing about negative value
	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestStartWithZeroIndexCurr(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(1)
	// tracker should do nothing about negative value
	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(1), tracker.indexCurr)
	assert.Equal(t, int64(1), tracker.lastTimestampMs)
	assert.Equal(t, int64(1), tracker.countTotal)
	// we don't track elapsed time in Start()
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestStartWithTwoIndexCurr(t *testing.T) {
	tracker := newTimeTracker("fake")

	tracker.Start(1)
	tracker.Start(2)

	// tracker should do nothing about negative value
	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(2), tracker.indexCurr)
	assert.Equal(t, int64(2), tracker.lastTimestampMs)
	assert.Equal(t, int64(2), tracker.countTotal)
	// we track elapsed time in Start() if called multiple times
	assert.Equal(t, int64(1), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestStartWithThreeIndexCurr(t *testing.T) {
	tracker := newTimeTracker("fake")

	tracker.Start(1)
	tracker.Start(2)
	tracker.Start(3)

	// tracker should do nothing about negative value
	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(3), tracker.indexCurr)
	assert.Equal(t, int64(3), tracker.lastTimestampMs)
	assert.Equal(t, int64(3), tracker.countTotal)
	// we don't track elapsed time in Start()
	assert.Equal(t, int64(3), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestEndWithoutStart(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.End(1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestEndWithNegativeNowMS(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.End(-1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestEndOneStart(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(1)
	tracker.End(2)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(2), tracker.lastTimestampMs)
	assert.Equal(t, int64(1), tracker.countTotal)
	assert.Equal(t, int64(1), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestEndTwoStart(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(1)
	tracker.End(2)

	tracker.Start(3)
	tracker.End(4)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(4), tracker.lastTimestampMs)
	assert.Equal(t, int64(2), tracker.countTotal)
	assert.Equal(t, int64(2), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestEndWithIncompleteEnd(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(1)

	tracker.Start(2)
	tracker.End(3)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(1), tracker.indexCurr)
	assert.Equal(t, int64(3), tracker.lastTimestampMs)
	assert.Equal(t, int64(2), tracker.countTotal)
	assert.Equal(t, int64(3), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestElapseWithNegativeParam(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Elapse(-1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestElapseHappyCase(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Elapse(1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(1), tracker.countTotal)
	assert.Equal(t, int64(1), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestElapseWithSampleWithNegativeTime(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.ElapseWithSample(-1, 1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestElapseWithSampleWithNegativeSample(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.ElapseWithSample(1, -1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestElapseWithSampleHappyCase(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.ElapseWithSample(1, 1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(1), tracker.countTotal)
	assert.Equal(t, int64(1), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestFinishHappyCase(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(1)
	tracker.End(2)

	tracker.Finish()

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(2), tracker.lastTimestampMs)
	assert.Equal(t, int64(1), tracker.countTotal)
	assert.Equal(t, int64(1), tracker.elapsedTotalMs)
	assert.True(t, tracker.isFinished)
}

func TestFinishWithoutEnd(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(1)

	tracker.Finish()

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(1), tracker.countTotal)
	assert.True(t, tracker.isFinished)
}
