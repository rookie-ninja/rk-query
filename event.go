// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import "go.uber.org/zap"

// An interface for recording event types.
// A typical event could be http GET/PUT/DELETE request.
type EventData interface {
	// Application name should be immutable
	GetApplicationName() string

	// Get hostname with system call
	GetHostName() string

	// Get event data name like GET / PUT / DELETE
	GetOperation() string

	// Set event data name like GET / PUT / DELETE
	SetOperation(string)

	// Reset all the event data fields
	Reset()

	GetEventStatus() eventDataStatus

	// Timer related
	GetStartTime() int64

	SetStartTime(int64) error

	GetEndTime() int64

	SetEndTime(endTime int64) error

	StartTimer(string) error

	EndTimer(string) error

	UpdateTimer(string, int64) error

	UpdateTimerWithSample(string, int64, int64) error

	// If no timer exists for the name, -1 is returned.
	GetTimeElapsed(string) int64

	// Remote Address related
	GetRemoteAddr() string

	SetRemoteAddr(string)

	// Counter related
	GetCounter(string) int64

	SetCounter(string, int64)
	// Increments the named counter by the given value.
	InCCounter(string, int64)

	AddKeyValuePair(key, value string)

	// appends the value to an existing name value pair.
	// If the value does not exist in the collection of
	// name value pairs, it is added.  If the value already
	// exists, a comma is appended to the existing value
	// and then the new value is appended.
	//
	// If appends more than DEFAULT_MAX_VALUE_COUNT
	// records, oldest records will be thrown out.
	AppendKeyValuePair(key, value string)

	GetValue(key string) string

	// Finish Event Data records the time of the last event.
	// Event time is measured as the time difference between calls to this method.
	// There are uses where we want to measure the time that we are in mutually exclusive states.
	// There are situations where determining the next or the previous state is
	// not possible at the same instant.  This method alleviates the problem.
	// You need not know want timer will be starting.
	// You just have to supply the timer that finished.
	//
	// One property of this method is that the sum of all events recorded
	// with this method will equal the sum the time for the entire event data.
	// (As long as start time is not changed)
	FinishCurrentEvent(name string) error

	// Inserts an event into the event history.
	RecordHistoryEvent(name string) error

	// Writes the event data to the query log.
	RecordProfiledData() error

	// Output as JSON format which is default
	ToJsonFormat() string

	ToZapFields() []zap.Field

	// Output as human friendly format
	ToPrettyFormat() string

	GetEventHistory() *eventHistory

	GetLogger() *zap.Logger
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

