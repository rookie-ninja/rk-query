// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkquery

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEventHistory(t *testing.T) {
	his := newEventHistory()

	assert.NotNil(t, his)
	assert.NotNil(t, his.builder)
	assert.False(t, his.truncated)
	assert.Zero(t, his.previousTimeMS)
}

func TestElapsedMS_WithTruncated(t *testing.T) {
	his := newEventHistory()
	his.truncated = true

	his.elapsedMS("action", 0)

	assert.Zero(t, his.builder.Len())
}

func TestElapsedMS(t *testing.T) {
	his := newEventHistory()
	his.elapsedMS("action", 1)
	assert.Equal(t, "action:1", his.builder.String())
}

func TestClear(t *testing.T) {
	his := newEventHistory()
	his.elapsedMS("action", 1)

	his.clear()

	assert.NotNil(t, his.builder)
	assert.False(t, his.truncated)
	assert.Zero(t, his.previousTimeMS)
}

func TestAppendTo(t *testing.T) {
	his := newEventHistory()
	his.elapsedMS("action", 1)

	buffer := &bytes.Buffer{}

	his.appendTo(buffer)
	assert.Equal(t, "action:1", buffer.String())
}

func TestGetHistoryLength(t *testing.T) {
	his := newEventHistory()
	his.elapsedMS("action", 1)

	assert.Equal(t, his.builder.Len(), his.getHistoryLength())
}

func TestEventHistory_String(t *testing.T) {
	his := newEventHistory()
	his.elapsedMS("action", 1)

	assert.Equal(t, his.builder.String(), his.String())
}
