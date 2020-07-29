// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEventZapFactoryHappyCase(t *testing.T) {
	fac := NewEventZapFactory()
	assert.NotNil(t, fac)
}

func TestEventZapFactory_CreateEvent(t *testing.T) {
	fac := NewEventZapFactory()
	event := fac.CreateEventZap()
	assert.NotNil(t, event)
}

func TestEventZapFactory_CreateEvent_WithAppName(t *testing.T) {
	fac := NewEventZapFactory()
	event := fac.CreateEventZap(WithAppName("app"))
	assert.NotNil(t, event)
	assert.Equal(t, "appName", event.GetAppName())
}
