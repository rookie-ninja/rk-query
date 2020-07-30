// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEventFactoryHappyCase(t *testing.T) {
	fac := NewEventFactory()
	assert.NotNil(t, fac)
}

func TestEventFactory_CreateEvent(t *testing.T) {
	fac := NewEventFactory()
	event := fac.CreateEvent()
	assert.NotNil(t, event)
}

func TestEventFactory_CreateEvent_WithAppName(t *testing.T) {
	fac := NewEventFactory()
	event := fac.CreateEvent(WithAppName("app"))
	assert.NotNil(t, event)
	assert.Equal(t, "appName", event.GetAppName())
}
