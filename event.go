package rkquery

import (
	"go.uber.org/zap"
	"time"
)

type Event interface {
	GetValue(string) string

	GetAppName() string

	GetEventId() string

	SetEventId(string)

	GetHostname() string

	GetLogger() *zap.Logger

	GetOperation() string

	SetOperation(string)

	GetEventStatus() eventStatus

	SetStartTime(time.Time)

	GetStartTime() time.Time

	GetEndTime() time.Time

	SetEndTime(time.Time)

	StartTimer(string)

	EndTimer(string)

	UpdateTimer(string, int64)

	UpdateTimerWithSample(string, int64, int64)

	GetTimeElapsedMS(string) int64

	GetRemoteAddr() string

	SetRemoteAddr(string)

	GetCounter(string) int64

	SetCounter(string, int64)

	InCCounter(string, int64)

	AddPair(string, string)

	AddErr(error)

	SetResCode(string)

	GetErrCount(error) int64

	AddFields(...zap.Field)

	GetFields() []zap.Field

	RecordHistoryEvent(string)

	WriteLog()

	setLogger(*zap.Logger)

	setFormat(format)

	setQuietMode(bool)

	setAppName(string)

	setHostname(string)
}
