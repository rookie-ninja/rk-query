// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"bytes"
	"github.com/juju/errors"
	"go.uber.org/zap"
	"strconv"
)

type timeTracker struct {
	name            string
	indexCurr       int64
	lastTimestampMS int64
	countTotal      int64
	elapsedTotalMS  int64
	isFinished      bool
}

func NewTimeTracker(name string) *timeTracker {
	if len(name) == 0 {
		return nil
	}

	return &timeTracker{
		name:            name,
		indexCurr:       0,
		lastTimestampMS: 0,
		countTotal:      0,
		elapsedTotalMS:  0,
		isFinished:      false,
	}
}

func (tracker *timeTracker) GetName() string {
	return tracker.name
}

func (tracker *timeTracker) GetCount() int64 {
	return tracker.countTotal
}

func (tracker *timeTracker) GetElapsedMS() int64 {
	return tracker.elapsedTotalMS
}

func (tracker *timeTracker) Start(nowMS int64) {
	if nowMS < 0 {
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
		tracker.elapsedTotalMS += tracker.indexCurr * (nowMS - tracker.lastTimestampMS)
	}

	tracker.lastTimestampMS = nowMS
	tracker.countTotal++
	tracker.indexCurr++
}

func (tracker *timeTracker) End(nowMS int64) {
	if tracker.indexCurr < 1 || nowMS < 0 {
		return
	}

	tracker.elapsedTotalMS += tracker.indexCurr * (nowMS - tracker.lastTimestampMS)
	tracker.lastTimestampMS = nowMS
	tracker.indexCurr--
}

func (tracker *timeTracker) Elapse(elapseTimeMS int64) {
	if elapseTimeMS < 0 {
		return
	}
	tracker.ElapseWithSample(elapseTimeMS, 1)
}

func (tracker *timeTracker) ElapseWithSample(elapseTimeMS int64, numSample int64) {
	if elapseTimeMS < 0 || numSample < 0 {
		return
	}

	tracker.countTotal += numSample
	tracker.elapsedTotalMS += elapseTimeMS
}

func (tracker *timeTracker) Finish(timeSource TimeSource) {
	tracker.isFinished = true

	if tracker.indexCurr == 0 {
		return
	}

	nowMS := timeSource.CurrentTimeMS()

	tracker.elapsedTotalMS += tracker.indexCurr * (nowMS - tracker.lastTimestampMS)
	tracker.lastTimestampMS = nowMS
	tracker.indexCurr = 0
}

func (tracker *timeTracker) String() string {
	var builder bytes.Buffer

	builder.WriteString(tracker.name + ":" +
		strconv.FormatInt(tracker.elapsedTotalMS, 10) + "/" +
		strconv.FormatInt(tracker.countTotal, 10))

	return builder.String()
}

func (tracker *timeTracker) StringWithTimeSource(timeSource TimeSource) string {
	if tracker.indexCurr == 0 {
		return tracker.String()
	}

	elapsed := tracker.elapsedTotalMS + tracker.indexCurr * (timeSource.CurrentTimeMS() - tracker.lastTimestampMS)

	var builder bytes.Buffer

	builder.WriteString(tracker.name + OpenMarker +
		strconv.FormatInt(tracker.indexCurr, 10) + ":" +
		strconv.FormatInt(elapsed, 10) + "/" +
		strconv.FormatInt(tracker.countTotal, 10))

	return builder.String()
}

func (tracker *timeTracker) ToZapFieldsWithTimeSource(timeSource TimeSource) ([]zap.Field, error) {
	if tracker.indexCurr == 0 {
		return tracker.ToZapFields()
	}

	elapsedMS := tracker.elapsedTotalMS + tracker.indexCurr * (timeSource.CurrentTimeMS() - tracker.lastTimestampMS)

	return []zap.Field{
		zap.Int64(tracker.name+OpenMarker + strconv.FormatInt(tracker.indexCurr, 10) + ".elapsed_ms", elapsedMS),
		zap.Int64(tracker.name+OpenMarker + strconv.FormatInt(tracker.indexCurr, 10) + ".count", tracker.countTotal),
	}, nil
}

func (tracker *timeTracker) ToZapFields() ([]zap.Field, error) {
	if tracker.indexCurr > 0 {
		return nil, errors.NewNotValid(nil,"there is still open timer")
	}

	return []zap.Field{
		zap.Int64(tracker.name+".elapsed_ms", tracker.elapsedTotalMS),
		zap.Int64(tracker.name+".count", tracker.countTotal),
	}, nil
}

