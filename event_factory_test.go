// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkquery

import (
	rklogger "github.com/rookie-ninja/rk-logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestWithZapLogger_WithNilLogger(t *testing.T) {
	opt := WithZapLogger(nil)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.NotNil(t, event.(*eventZap).logger)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.NotNil(t, threadSafe.(*eventThreadSafe).delegate.logger)
}

func TestWithZapLogger_HappyCase(t *testing.T) {
	logger := rklogger.NoopLogger
	opt := WithZapLogger(logger)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.Equal(t, logger, event.(*eventZap).logger)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.Equal(t, logger, threadSafe.(*eventThreadSafe).delegate.logger)
}

func TestWithEncoding_InvalidEncodingType(t *testing.T) {
	opt := WithEncoding(-1)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.NotEqual(t, -1, event.(*eventZap).encoding)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.NotEqual(t, -1, threadSafe.(*eventThreadSafe).delegate.encoding)
}

func TestWithEncoding_WithJson(t *testing.T) {
	opt := WithEncoding(JSON)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.Equal(t, JSON.String(), event.(*eventZap).encoding.String())

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.Equal(t, JSON.String(), threadSafe.(*eventThreadSafe).delegate.encoding.String())
}

func TestWithEncoding_WithConsole(t *testing.T) {
	opt := WithEncoding(CONSOLE)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.Equal(t, CONSOLE.String(), event.(*eventZap).encoding.String())

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.Equal(t, CONSOLE.String(), threadSafe.(*eventThreadSafe).delegate.encoding.String())
}

func TestWithQuietMode_WithTrue(t *testing.T) {
	opt := WithQuietMode(true)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.True(t, event.(*eventZap).quietMode)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.True(t, threadSafe.(*eventThreadSafe).delegate.quietMode)
}

func TestWithQuietMode_WithFalse(t *testing.T) {
	opt := WithQuietMode(false)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.False(t, event.(*eventZap).quietMode)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.False(t, threadSafe.(*eventThreadSafe).delegate.quietMode)
}

func TestWithEntryName_HappyCase(t *testing.T) {
	entryName := "ut-entry"
	opt := WithEntryName(entryName)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.Equal(t, entryName, event.(*eventZap).entryName)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.Equal(t, entryName, threadSafe.(*eventThreadSafe).delegate.entryName)
}

func TestWithEntryType_HappyCase(t *testing.T) {
	entryKind := "ut-kind"
	opt := WithEntryKind(entryKind)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.Equal(t, entryKind, event.(*eventZap).entryKind)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.Equal(t, entryKind, threadSafe.(*eventThreadSafe).delegate.entryKind)
}

func TestWithWithAppName_HappyCase(t *testing.T) {
	appName := "ut-service"
	opt := WithServiceName(appName)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.Equal(t, appName, event.(*eventZap).serviceName)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.Equal(t, appName, threadSafe.(*eventThreadSafe).delegate.serviceName)
}

func TestWithWithAppVersion_HappyCase(t *testing.T) {
	appVersion := "ut-version"
	opt := WithServiceVersion(appVersion)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.Equal(t, appVersion, event.(*eventZap).serviceVersion)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.Equal(t, appVersion, threadSafe.(*eventThreadSafe).delegate.serviceVersion)
}

func TestWithWithOperation_HappyCase(t *testing.T) {
	operation := "ut-operation"
	opt := WithOperation(operation)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.Equal(t, operation, event.(*eventZap).operation)

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.Equal(t, operation, threadSafe.(*eventThreadSafe).delegate.operation)
}

func TestWithWithPayloads_HappyCase(t *testing.T) {
	field := zap.String("key", "value")
	opt := WithPayloads([]zap.Field{field}...)

	// For eventZap
	event := NewEventFactory().CreateEvent()
	opt(event)
	assert.Len(t, event.(*eventZap).payloads, 1)
	assert.Equal(t, field, event.(*eventZap).payloads[0])

	// For eventThreadSafe
	threadSafe := NewEventFactory().CreateEventThreadSafe()
	opt(threadSafe)
	assert.Len(t, threadSafe.(*eventThreadSafe).delegate.payloads, 1)
	assert.Equal(t, field, threadSafe.(*eventThreadSafe).delegate.payloads[0])
}

func TestEventFactory_CreateEvent_WithoutOptions(t *testing.T) {
	fac := NewEventFactory()
	event := fac.CreateEvent()
	assert.NotNil(t, event)
	// Check default fields
	assert.Equal(t, rklogger.EventLogger, event.(*eventZap).logger)
	assert.Equal(t, CONSOLE.String(), event.(*eventZap).encoding.String())
	assert.False(t, event.(*eventZap).quietMode)
	assert.Empty(t, event.(*eventZap).serviceName)
	assert.Empty(t, event.(*eventZap).serviceVersion)
	assert.Empty(t, event.(*eventZap).entryName)
	assert.Empty(t, event.(*eventZap).entryKind)
	assert.NotEmpty(t, event.(*eventZap).eventId)
	assert.Empty(t, event.(*eventZap).traceId)
	assert.Empty(t, event.(*eventZap).requestId)
	assert.NotNil(t, event.(*eventZap).startTime)
	assert.NotEmpty(t, event.(*eventZap).timeZone)
	assert.Empty(t, event.(*eventZap).payloads)
	assert.NotNil(t, event.(*eventZap).errors)
	assert.Empty(t, event.(*eventZap).operation)
	assert.Equal(t, "localhost", event.(*eventZap).remoteAddr)
	assert.Empty(t, event.(*eventZap).resCode)
	assert.Equal(t, NotStarted, event.(*eventZap).status)
	assert.NotNil(t, event.(*eventZap).pairs)
	assert.NotNil(t, event.(*eventZap).counters)
	assert.NotNil(t, event.(*eventZap).tracker)
}

func TestEventFactory_CreateEvent_WithOptions(t *testing.T) {
	fac := NewEventFactory(WithServiceName("ut-service"))
	event := fac.CreateEvent(WithServiceVersion("ut-version"))
	assert.NotNil(t, event)

	// Check default fields
	assert.Equal(t, "ut-service", event.(*eventZap).serviceName)
	assert.Equal(t, "ut-version", event.(*eventZap).serviceVersion)
}

func TestNewEventFactory_HappyCase(t *testing.T) {
	fac := NewEventFactory()
	assert.NotNil(t, fac)
}

func TestEventFactory_CreateEventNoop(t *testing.T) {
	fac := NewEventFactory()
	assert.NotNil(t, fac.CreateEventNoop())
}

func TestEventFactory_CreateEventThreadSafe(t *testing.T) {
	fac := NewEventFactory()
	event := fac.CreateEventThreadSafe()
	assert.NotNil(t, event)
	assert.NotNil(t, event.(*eventThreadSafe).lock)
	assert.NotNil(t, event.(*eventThreadSafe).delegate)
}

func TestGetDefaultIfEmptyString_WithEmptyOrigin(t *testing.T) {
	res := getDefaultIfEmptyString("", "ut-default")
	assert.Equal(t, "ut-default", res)
}

func TestGetDefaultIfEmptyString_WithNonEmptyOrigin(t *testing.T) {
	res := getDefaultIfEmptyString("ut-origin", "ut-default")
	assert.Equal(t, "ut-origin", res)
}
