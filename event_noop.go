package rk_query

import (
	"go.uber.org/zap"
	"time"
)

type EventNoop struct{}

func (EventNoop) GetValue(string) string {
	return ""
}

func (EventNoop) GetAppName() string {
	return ""
}

func (EventNoop) GetEventId() string {
	return ""
}

func (EventNoop) SetEventId(string) {}

func (EventNoop) GetHostname() string {
	return ""
}

func (EventNoop) GetLogger() *zap.Logger {
	return zap.NewNop()
}

func (EventNoop) GetOperation() string {
	return ""
}

func (EventNoop) SetOperation(string) {}

func (EventNoop) GetEventStatus() eventStatus {
	return notStarted
}

func (EventNoop) SetStartTime(time.Time) {}

func (EventNoop) GetStartTime() time.Time {
	return time.Now()
}

func (EventNoop) GetEndTime() time.Time {
	return time.Now()
}

func (EventNoop) SetEndTime(time.Time) {}

func (EventNoop) StartTimer(string) {}

func (EventNoop) EndTimer(string) {}

func (EventNoop) UpdateTimer(string, int64) {}

func (EventNoop) UpdateTimerWithSample(string, int64, int64) {}

func (EventNoop) GetTimeElapsedMS(string) int64 {
	return 0
}

func (EventNoop) GetRemoteAddr() string {
	return ""
}

func (EventNoop) SetRemoteAddr(string) {}

func (EventNoop) GetCounter(string) int64 {
	return 0
}

func (EventNoop) SetCounter(string, int64) {}

func (EventNoop) InCCounter(string, int64) {}

func (EventNoop) AddPair(string, string) {}

func (EventNoop) AddErr(error) {
}

func (EventNoop) GetErrCount(error) int64 {
	return 0
}

func (EventNoop) AddFields(...zap.Field) {}

func (EventNoop) GetFields() []zap.Field {
	return make([]zap.Field, 0)
}

func (EventNoop) RecordHistoryEvent(string) {}

func (EventNoop) WriteLog() {}

func (EventNoop) setLogger(*zap.Logger) {}

func (EventNoop) setFormat(Format) {}

func (EventNoop) setQuietMode(bool) {}

func (EventNoop) setAppName(string) {}

func (EventNoop) setHostname(string) {}
