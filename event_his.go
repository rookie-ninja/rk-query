// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"bytes"
	"strconv"
)

type eventHistory struct {
	builder        bytes.Buffer
	truncated      bool
	previousTimeMS int64
}

func newEventHistory() *eventHistory {
	return &eventHistory{
		builder:        bytes.Buffer{},
		truncated:      false,
		previousTimeMS: 0,
	}
}

// Elapse with action in milliseconds
func (his *eventHistory) elapsedMS(action string, elapseMS int64) {
	if his.truncated {
		return
	}

	// do we have enough space?
	length := his.builder.Len()
	elapsed := strconv.FormatInt(elapseMS-his.previousTimeMS, 10)

	size := len(action) + 1 + len(elapsed)
	if length > 0 {
		size++
	}

	if length+size+1+len(truncated) > maxHistoryLength {
		his.truncated = true
		if length > 0 {
			// we have something in the string and adding more would've
			// put us over our limit, so just mark the string truncated
			his.builder.WriteString(commaTruncated)
		} else {
			// we have nothing in the string and we were asked to add
			// something so large that we'd immediately be over the limit;
			// we'll immediately go TRUNCATED in this case.
			his.builder.WriteString(truncated)
		}
		return
	}

	// save the previous time
	his.previousTimeMS = elapseMS

	// append a comma if we've got a previous string
	if length > 1 {
		his.builder.WriteByte(',')
	}

	// append our next action, a colon, and its delta time
	his.builder.WriteString(action)
	his.builder.WriteByte(':')
	his.builder.WriteString(elapsed)
}

// Clear
func (his *eventHistory) clear() {
	his.builder.Reset()
	his.previousTimeMS = 0
	his.truncated = false
}

// Appends the history to given buffer
func (his *eventHistory) appendTo(buffer *bytes.Buffer) {
	buffer.Write(his.builder.Bytes())
}

// Length of current history
func (his *eventHistory) getHistoryLength() int {
	return his.builder.Len()
}

func (his *eventHistory) String() string {
	return his.builder.String()
}
