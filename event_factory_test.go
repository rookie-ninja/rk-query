// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestNewEventFactoryWithNilZapLogger(t *testing.T) {
	fac := NewEventFactory("", &RealTimeSource{}, nil)
	assert.NotNil(t, fac)
}

func TestNewEventFactoryHappyCase(t *testing.T) {
	fac := NewEventFactory("", &RealTimeSource{}, zap.NewNop())
	assert.NotNil(t, fac)
	assert.NotNil(t, fac.TimeSource)
	assert.Empty(t, fac.AppName)
	assert.NotEmpty(t, fac.HostName)
	assert.NotNil(t, fac.ZapLogger)
	assert.Equal(t, RK, fac.Format)
	assert.False(t, fac.Minimal)
	assert.Empty(t, fac.Listeners)
	assert.Empty(t, fac.DefaultKvs)
}

func TestEventFactory_CreateEvent(t *testing.T) {
	fac := NewEventFactory("", &RealTimeSource{}, zap.NewNop())
	event := fac.CreateEvent()
	assert.NotNil(t, event)
	assert.IsType(t, &EventImpl{}, event)
}

func TestEventFactory_CreateThreadSafeEvent(t *testing.T) {
	fac := NewEventFactory("", &RealTimeSource{}, zap.NewNop())
	event := fac.CreateThreadSafeEvent()
	assert.NotNil(t, event)
	assert.IsType(t, &ThreadSafeEventImpl{}, event)
}

func TestEventFactory_CreateNoopEvent(t *testing.T) {
	fac := NewEventFactory("", &RealTimeSource{}, zap.NewNop())
	event := fac.CreateNoopEvent()
	assert.NotNil(t, event)
	assert.IsType(t, &NoopEvent{}, event)
}
