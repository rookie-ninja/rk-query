package rk_query

import (
	"go.uber.org/zap"
	"time"
)

type eventNoop struct{}

func (eventNoop) GetValue(string) string {
	return ""
}

func (eventNoop) GetAppName() string {
	return ""
}

func (eventNoop) GetEventId() string {
	return ""
}

func (eventNoop) SetEventId(string) {}

func (eventNoop) GetHostname() string {
	return ""
}

func (eventNoop) GetLogger() *zap.Logger {
	return zap.NewNop()
}

func (eventNoop) GetOperation() string {
	return ""
}

func (eventNoop) SetOperation(string) {}

func (eventNoop) GetEventStatus() eventStatus {
	return notStarted
}

func (eventNoop) SetStartTime(time.Time) {}

func (eventNoop) GetStartTime() time.Time {
	return time.Now()
}

func (eventNoop) GetEndTime() time.Time {
	return time.Now()
}

func (eventNoop) SetEndTime(time.Time) {}

func (eventNoop) StartTimer(string) {}

func (eventNoop) EndTimer(string) {}

func (eventNoop) UpdateTimer(string, int64) {}

func (eventNoop) UpdateTimerWithSample(string, int64, int64) {}

func (eventNoop) GetTimeElapsedMS(string) int64 {
	return 0
}

func (eventNoop) GetRemoteAddr() string {
	return ""
}

func (eventNoop) SetRemoteAddr(string) {}

func (eventNoop) GetCounter(string) int64 {
	return 0
}

func (eventNoop) SetCounter(string, int64) {}

func (eventNoop) InCCounter(string, int64) {}

func (eventNoop) AddPair(string, string) {}

func (eventNoop) AddErr(error) {
}

func (eventNoop) GetErrCount(error) int64 {
	return 0
}

func (eventNoop) AddFields(...zap.Field) {}

func (eventNoop) GetFields() []zap.Field {
	return make([]zap.Field, 0)
}

func (eventNoop) RecordHistoryEvent(string) {}

func (eventNoop) WriteLog() {}

func (eventNoop) setLogger(*zap.Logger) {}

func (eventNoop) setFormat(format) {}

func (eventNoop) setQuietMode(bool) {}

func (eventNoop) setAppName(string) {}

func (eventNoop) setHostname(string) {}
