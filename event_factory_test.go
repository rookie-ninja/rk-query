// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package rkquery

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
	assert.Equal(t, "app", event.(*eventZap).appName)
}
