// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

type timeTracker struct {
	name            string
	indexCurr       int64
	lastTimestampMs int64
	countTotal      int64
	elapsedTotalMs  int64
	isFinished      bool
}

// Create a new timeTracker with name.
// Name should be unique.
func newTimeTracker(name string) *timeTracker {
	if len(name) == 0 {
		return nil
	}

	return &timeTracker{
		name:            name,
		indexCurr:       0,
		lastTimestampMs: 0,
		countTotal:      0,
		elapsedTotalMs:  0,
		isFinished:      false,
	}
}

// Get name of current timeTracker.
func (tracker *timeTracker) GetName() string {
	return tracker.name
}

// Get count of how many times Start() has been called.
func (tracker *timeTracker) GetCount() int64 {
	return tracker.countTotal
}

// Get elapsed time in milli seconds.
func (tracker *timeTracker) GetElapsedMs() int64 {
	return tracker.elapsedTotalMs
}

// Start current timer with provided timestamp.
func (tracker *timeTracker) Start(nowMs int64) {
	if nowMs < 0 {
		return
	}

	// This is a little bit hard to understand bellow logic
	// For example, for every duplicated event, we track every calls.
	//
	// |  22  333
	// |111111111111
	// +-------------
	//
	// In case above, the x-axis represents the time, y-axis represents concurrent event.
	// The total time elapsed would be 18 with this formula: 17 = 12 + 3 + 2
	if tracker.indexCurr > 0 {
		tracker.elapsedTotalMs += tracker.indexCurr * (nowMs - tracker.lastTimestampMs)
	}

	tracker.lastTimestampMs = nowMs
	tracker.countTotal++
	tracker.indexCurr++
}

// End current timer with provided timestamp.
func (tracker *timeTracker) End(nowMs int64) {
	if tracker.indexCurr < 1 || nowMs < 0 {
		return
	}

	tracker.elapsedTotalMs += tracker.indexCurr * (nowMs - tracker.lastTimestampMs)
	tracker.lastTimestampMs = nowMs
	tracker.indexCurr--
}

// Force to elapse timer.
func (tracker *timeTracker) Elapse(elapseTimeMs int64) {
	if elapseTimeMs < 0 {
		return
	}
	tracker.ElapseWithSample(elapseTimeMs, 1)
}

// For to elapse timer with number of sample.
func (tracker *timeTracker) ElapseWithSample(elapseTimeMs int64, numSample int64) {
	if elapseTimeMs < 0 || numSample < 0 {
		return
	}

	tracker.countTotal += numSample
	tracker.elapsedTotalMs += elapseTimeMs
}

// stop current timer.
func (tracker *timeTracker) Finish() {
	tracker.isFinished = true

	if tracker.indexCurr == 0 {
		return
	}

	nowMs := toMillisecond(time.Now())

	tracker.elapsedTotalMs += tracker.indexCurr * (nowMs - tracker.lastTimestampMs)
	tracker.lastTimestampMs = nowMs
	tracker.indexCurr = 0
}

// Convert to zap fields.
func (tracker *timeTracker) ToZapFields(enc *zapcore.MapObjectEncoder) []zap.Field {
	if tracker.indexCurr == 0 {
		if enc != nil {
			enc.AddInt64(tracker.name+".elapsedMs", tracker.elapsedTotalMs)
			enc.AddInt64(tracker.name+".count", tracker.countTotal)
		}

		return []zap.Field{
			zap.Int64(tracker.name+".elapsedMs", tracker.elapsedTotalMs),
			zap.Int64(tracker.name+".count", tracker.countTotal),
		}
	}

	nowMs := toMillisecond(time.Now())
	elapsedMs := tracker.elapsedTotalMs + tracker.indexCurr*(nowMs-tracker.lastTimestampMs)

	if enc != nil {
		enc.AddInt64(tracker.name+openMarker+strconv.FormatInt(tracker.indexCurr, 10)+".elapsedMs", tracker.elapsedTotalMs)
		enc.AddInt64(tracker.name+openMarker+strconv.FormatInt(tracker.indexCurr, 10)+".count", tracker.countTotal)
	}
	return []zap.Field{
		zap.Int64(tracker.name+openMarker+strconv.FormatInt(tracker.indexCurr, 10)+".elapsedMs", elapsedMs),
		zap.Int64(tracker.name+openMarker+strconv.FormatInt(tracker.indexCurr, 10)+".count", tracker.countTotal),
	}
}
