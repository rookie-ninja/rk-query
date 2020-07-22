// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import "go.uber.org/zap"

// An interface for recording event types.
// A typical event could be http GET/PUT/DELETE request.
type Event interface {
	// Application name should be immutable
	GetAppName() string

	// Get hostname with system call
	GetHostName() string

	// Get zap logger underlying
	GetZapLogger() *zap.Logger

	// Get event data name like GET / PUT / DELETE
	GetOperation() string

	// Set event data name like GET / PUT / DELETE
	SetOperation(string)

	// Reset all the event data fields
	Reset()

	GetEventStatus() eventDataStatus

	GetEndTimeMS() int64

	GetStartTimeMS() int64

	SetStartTimeMS(int64)

	SetEndTimeMS(int64)

	StartTimer(string)

	EndTimer(string)

	UpdateTimer(string, int64)

	UpdateTimerWithSample(string, int64, int64)

	// If no timer exists for the name, -1 is returned.
	GetTimeElapsedMS(string) int64

	// Remote Address related
	GetRemoteAddr() string

	SetRemoteAddr(string)

	// Counter related
	GetCounter(string) int64

	SetCounter(string, int64)
	// Increments the named counter by the given value.
	InCCounter(string, int64)

	AddErr(error)

	AddKv(string, string)

	AppendKv(string, string)

	GetValue(string) string

	FinishCurrentEvent(string)

	// Inserts an event into the event history.
	RecordHistoryEvent(string)

	// Writes the event data to log.
	WriteLog()

	ToZapFieldsMin() []zap.Field

	ToZapFields() []zap.Field

	GetEventHistory() *eventHistory
}

type eventDataStatus int

const (
	notStarted eventDataStatus = 0
	inProgress eventDataStatus = 1
	ended      eventDataStatus = 2
)

func (status eventDataStatus) String() string {
	names := [...]string{
		"NotStarted",
		"InProgress",
		"Ended",
	}

	if status < notStarted || status > ended {
		return "Unknown"
	}

	return names[status]
}

type Format int

const (
	JSON Format = 0
	RK   Format = 1
)

// Stringfy above config file types.
func (fileType Format) String() string {
	names := [...]string{"JSON", "RK"}

	// Please do not forget to change the boundary while adding a new config file types
	if fileType < JSON || fileType > RK {
		return "UNKNOWN"
	}

	return names[fileType]
}
