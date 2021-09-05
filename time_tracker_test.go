// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkquery

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestNewTimeTracker_WithNilName(t *testing.T) {
	assert.Nil(t, newTimeTracker(""))
}

func TestNewTimeTracker_HappyCase(t *testing.T) {
	tracker := newTimeTracker("fake")
	assert.NotNil(t, tracker)
	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestTimeTracker_GetName(t *testing.T) {
	assert.Equal(t, "fake", newTimeTracker("fake").GetName())
}

func TestTimeTracker_GetCount(t *testing.T) {
	assert.Zero(t, newTimeTracker("fake").GetCount())
}

func TestTimeTracker_GetElapsedMs(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.elapsedTotalMs = 1
	assert.Equal(t, tracker.elapsedTotalMs, tracker.GetElapsedMs())
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

func TestStart_WithZeroIndexCurr(t *testing.T) {
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

func TestStart_WithTwoIndexCurr(t *testing.T) {
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

func TestStart_WithThreeIndexCurr(t *testing.T) {
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

func TestEnd_WithoutStart(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.End(1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestEnd_WithNegativeNowMS(t *testing.T) {
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

func TestElapse_WithNegativeParam(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Elapse(-1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestElapse_HappyCase(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Elapse(1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(1), tracker.countTotal)
	assert.Equal(t, int64(1), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestElapseWithSample_WithNegativeTime(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.ElapseWithSample(-1, 1)

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(0), tracker.lastTimestampMs)
	assert.Equal(t, int64(0), tracker.countTotal)
	assert.Equal(t, int64(0), tracker.elapsedTotalMs)
	assert.False(t, tracker.isFinished)
}

func TestElapseWithSample_WithNegativeSample(t *testing.T) {
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

func TestFinish_HappyCase(t *testing.T) {
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

func TestFinish_WithoutEnd(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(1)

	tracker.Finish()

	assert.Equal(t, "fake", tracker.name)
	assert.Equal(t, int64(0), tracker.indexCurr)
	assert.Equal(t, int64(1), tracker.countTotal)
	assert.True(t, tracker.isFinished)
}

func TestTimeTracker_ToZapFields_WithNilEncoder(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(time.Now().UnixNano())
	tracker.Elapse(1)
	tracker.Finish()

	assert.NotEmpty(t, tracker.ToZapFields(nil))
}

func TestTimeTracker_ToZapFields_HappyCase(t *testing.T) {
	tracker := newTimeTracker("fake")
	tracker.Start(time.Now().UnixNano())
	tracker.Elapse(1)
	tracker.Finish()

	assert.NotEmpty(t, tracker.ToZapFields(zapcore.NewMapObjectEncoder()))
}
