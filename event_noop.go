package rk_query

import "go.uber.org/zap"

type NoopEventData struct{}

func (NoopEventData) GetAppName() string {
	return ""
}

func (NoopEventData) GetHostName() string {
	return ""
}

func (NoopEventData) GetZapLogger() *zap.Logger {
	return nil
}

func (NoopEventData) GetOperation() string {
	return ""
}

func (NoopEventData) SetOperation(string) {}

func (NoopEventData) Reset() {}

func (NoopEventData) GetEventStatus() eventDataStatus {
	return 0
}

func (NoopEventData) GetEndTimeMS() int64 {
	return 0
}

func (NoopEventData) GetStartTimeMS() int64 {
	return 0
}

func (NoopEventData) SetStartTimeMS(int64) {}

func (NoopEventData) SetEndTimeMS(int64) {}

func (NoopEventData) StartTimer(string) {}

func (NoopEventData) EndTimer(string) {}

func (NoopEventData) UpdateTimer(string, int64) {}

func (NoopEventData) UpdateTimerWithSample(string, int64, int64) {}

func (NoopEventData) GetTimeElapsedMS(string) int64 {
	return 0
}

func (NoopEventData) GetRemoteAddr() string {
	return ""
}

func (NoopEventData) SetRemoteAddr(string) {}

func (NoopEventData) GetCounter(string) int64 {
	return 0
}

func (NoopEventData) SetCounter(string, int64) {}

func (NoopEventData) InCCounter(string, int64) {}

func (NoopEventData) AddKv(string, string) {}

func (NoopEventData) AddErr(error) {}

func (NoopEventData) AppendKv(string, string) {}

func (NoopEventData) GetValue(string) string {
	return ""
}

func (NoopEventData) FinishCurrentEvent(string) {}

func (NoopEventData) RecordHistoryEvent(string) {}

func (NoopEventData) WriteLog() {}

func (NoopEventData) ToZapFieldsMin() []zap.Field {
	return make([]zap.Field, 0)
}

func (NoopEventData) ToZapFields() []zap.Field {
	return make([]zap.Field, 0)
}

func (NoopEventData) GetEventHistory() *eventHistory {
	return nil
}