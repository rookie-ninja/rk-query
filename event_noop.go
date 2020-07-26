// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import "go.uber.org/zap"

type NoopEvent struct{}

func (NoopEvent) GetAppName() string {
	return ""
}

func (NoopEvent) GetHostName() string {
	return ""
}

func (NoopEvent) GetZapLogger() *zap.Logger {
	return nil
}

func (NoopEvent) GetOperation() string {
	return ""
}

func (NoopEvent) SetOperation(string) {}

func (NoopEvent) Reset() {}

func (NoopEvent) GetEventStatus() eventStatus {
	return 0
}

func (NoopEvent) GetEndTimeMS() int64 {
	return 0
}

func (NoopEvent) GetStartTimeMS() int64 {
	return 0
}

func (NoopEvent) SetStartTimeMS(int64) {}

func (NoopEvent) SetEndTimeMS(int64) {}

func (NoopEvent) StartTimer(string) {}

func (NoopEvent) EndTimer(string) {}

func (NoopEvent) UpdateTimer(string, int64) {}

func (NoopEvent) UpdateTimerWithSample(string, int64, int64) {}

func (NoopEvent) GetTimeElapsedMS(string) int64 {
	return 0
}

func (NoopEvent) GetRemoteAddr() string {
	return ""
}

func (NoopEvent) SetRemoteAddr(string) {}

func (NoopEvent) GetCounter(string) int64 {
	return 0
}

func (NoopEvent) SetCounter(string, int64) {}

func (NoopEvent) InCCounter(string, int64) {}

func (NoopEvent) AddKv(string, string) {}

func (NoopEvent) AddErr(error) {}

func (NoopEvent) GetErrCount(error) int64 {
	return 0
}

func (NoopEvent) AppendKv(string, string) {}

func (NoopEvent) GetValue(string) string {
	return ""
}

func (NoopEvent) FinishCurrentTimer(string) {}

func (NoopEvent) RecordHistoryEvent(string) {}

func (NoopEvent) WriteLog() {}

func (NoopEvent) ToZapFieldsMin() []zap.Field {
	return make([]zap.Field, 0)
}

func (NoopEvent) ToZapFields() []zap.Field {
	return make([]zap.Field, 0)
}

func (NoopEvent) GetEventHistory() *eventHistory {
	return nil
}
