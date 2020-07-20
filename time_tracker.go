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

	tracker := timeTracker{
		name:            name,
		indexCurr:       0,
		lastTimestampMS: 0,
		countTotal:      0,
		elapsedTotalMS:  0,
		isFinished:      false,
	}

	return &tracker
}

func (tracker *timeTracker) GetName() string {
	return tracker.name
}

func (tracker *timeTracker) Start(nowMS int64) {
	if tracker.indexCurr > 0 {
		tracker.elapsedTotalMS += tracker.indexCurr * (nowMS - tracker.lastTimestampMS)
	}

	tracker.lastTimestampMS = nowMS
	tracker.countTotal++
	tracker.indexCurr++
}

func (tracker *timeTracker) End(nowMS int64) {
	if tracker.indexCurr < 1 {
		return
	}

	tracker.elapsedTotalMS += tracker.indexCurr * (nowMS - tracker.lastTimestampMS)
	tracker.lastTimestampMS = nowMS
	tracker.indexCurr--
}

func (tracker *timeTracker) Elapse(elapseTimeMS int64) {
	tracker.ElapseWithSample(elapseTimeMS, 1)
}

func (tracker *timeTracker) ElapseWithSample(elapseTimeMS int64, numSample int64) {
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

func (tracker *timeTracker) StringWithTimeSource(timeSource TimeSource) (string, error) {
	if tracker.indexCurr == 0 {
		return tracker.String()
	}

	elapsed := tracker.elapsedTotalMS + tracker.indexCurr * (timeSource.CurrentTimeMS() - tracker.lastTimestampMS)

	var builder bytes.Buffer

	builder.WriteString(tracker.name + EventOpenMarker +
		strconv.FormatInt(tracker.indexCurr, 10) + ":" +
		strconv.FormatInt(elapsed, 10) + "/" +
		strconv.FormatInt(tracker.countTotal, 10))

	return builder.String(), nil
}

func (tracker *timeTracker) ToZapFieldsWithTimeSource(timeSource TimeSource) ([]zap.Field, error) {
	if tracker.indexCurr == 0 {
		return tracker.ToZapFields()
	}

	elapsedMS := tracker.elapsedTotalMS + tracker.indexCurr * (timeSource.CurrentTimeMS() - tracker.lastTimestampMS)

	return []zap.Field{
		zap.Int64(tracker.name+EventOpenMarker + strconv.FormatInt(tracker.indexCurr, 10) + ".elapsed_ms", elapsedMS),
		zap.Int64(tracker.name+EventOpenMarker + strconv.FormatInt(tracker.indexCurr, 10) + ".count", tracker.countTotal),
	}, nil
}

func (tracker *timeTracker) ToZapFields() ([]zap.Field, error) {
	if tracker.indexCurr > 0 {
		return nil, errors.New("there is still open timer")
	}

	return []zap.Field{
		zap.Int64(tracker.name+".elapsed_ms", tracker.elapsedTotalMS),
		zap.Int64(tracker.name+".count", tracker.countTotal),
	}, nil
}

func (tracker *timeTracker) String() (string, error) {
	if tracker.indexCurr > 0 {
		return "", errors.New("cannot call ToString() with open timers")
	}

	var builder bytes.Buffer

	builder.WriteString(tracker.name + ":" +
		strconv.FormatInt(tracker.elapsedTotalMS, 10) + "/" +
		strconv.FormatInt(tracker.countTotal, 10))

	return builder.String(), nil
}

func (tracker *timeTracker) GetCount() int64 {
	return tracker.countTotal
}

func (tracker *timeTracker) GetElapsedMS() int64 {
	return tracker.elapsedTotalMS
}
