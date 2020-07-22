// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"bytes"
	"strconv"
)

type eventHistory struct {
	builder        bytes.Buffer
	truncated      bool
	previousMillis int64
}

func newEventHistory() *eventHistory {
	his := eventHistory{}

	his.builder = bytes.Buffer{}
	his.truncated = false
	his.previousMillis = 0
	return &his
}

// Appends an elapsed time entry to the history.
func (his *eventHistory) elapsedMS(action string, time int64) {
	if his.truncated {
		return
	}

	// see if we have space for more history.
	length := his.builder.Len()
	elapsed := strconv.FormatInt(time-his.previousMillis, 10)

	size := len(action) + 1 + len(elapsed)
	if length > 0 {
		size++
	}

	if length+size+1+len(Truncated) > MaxHistoryLength {
		his.truncated = true
		if length > 0 {
			// we have something in the string and adding more would've
			// put us over our limit, so just mark the string truncated
			his.builder.WriteString(CommaTruncated)
		} else {
			// we have nothing in the string and we were asked to add
			// something so large that we'd immediately be over the limit;
			// we'll immediately go TRUNCATED in this case.
			his.builder.WriteString(Truncated)
		}
		return
	}

	// save the previous time
	his.previousMillis = time

	// append a comma if we've got a previous string
	if length > 1 {
		his.builder.WriteByte(',')
	}

	// append our next action, a colon, and its delta time
	his.builder.WriteString(action)
	his.builder.WriteByte(':')
	his.builder.WriteString(elapsed)
}

// Tearing down object
func (his *eventHistory) clear() {
	his.builder.Reset()
	his.previousMillis = 0
	his.truncated = false
}

// Appends the event history data to the given byte buffer
func (his *eventHistory) appendTo(buffer *bytes.Buffer) {
	buffer.Write(his.builder.Bytes())
}

// The length of strings
func (his *eventHistory) getHistoryLength() int {
	return his.builder.Len()
}

// Returns the current history as string
func (his *eventHistory) String() string {
	return his.builder.String()
}
