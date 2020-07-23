// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestFormatAsRkHappyCase(t *testing.T) {
	entry := &eventEntryImpl{}
	res := entry.FormatAsRk(getEvent())

	assert.True(t, strings.HasPrefix(res, ScopeDelimiter))
	assert.True(t, strings.HasSuffix(res, EOE))
	assert.True(t, strings.Contains(res, "end_time="))
	assert.True(t, strings.Contains(res, "start_time="))
	assert.True(t, strings.Contains(res, "time="))
	assert.True(t, strings.Contains(res, "hostname="))
	assert.True(t, strings.Contains(res, "app="))
	assert.True(t, strings.Contains(res, "operation="))
	assert.True(t, strings.Contains(res, "event_status="))
}

func TestFormatAsRkMinHappyCase(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent().(*EventImpl)

	res := entry.FormatAsRkMin(event)

	assert.True(t, strings.HasPrefix(res, ScopeDelimiter))
	assert.True(t, strings.HasSuffix(res, EOE))
	assert.False(t, strings.Contains(res, "end_time="))
	assert.False(t, strings.Contains(res, "start_time="))
	assert.False(t, strings.Contains(res, "time="))
	assert.False(t, strings.Contains(res, "hostname="))
	assert.False(t, strings.Contains(res, "app="))
	assert.False(t, strings.Contains(res, "operation="))
	assert.False(t, strings.Contains(res, "event_status="))
}

func TestGetKvAsZapFieldsWithEmptyKvs(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent()

	assert.Empty(t, entry.getCounterAsZapFields(event.(*EventImpl)))
}

func TestGetKvAsZapFieldsHappyCase(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent().(*EventImpl)
	event.AddKv("key", "value")

	assert.NotEmpty(t, entry.getKvAsZapFields(event))
	assert.Len(t, entry.getKvAsZapFields(event), 1)
	assert.Equal(t, "key", entry.getKvAsZapFields(event)[0].Key)
}

func TestGetErrAsZapFieldsWithEmptyErrs(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent().(*EventImpl)

	assert.Empty(t, entry.getErrAsZapFields(event))
}

func TestGetErrAsZapFieldsHappyCase(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent().(*EventImpl)
	event.AddErr(MyErr{})

	assert.NotEmpty(t, entry.getErrAsZapFields(event))
	assert.Len(t, entry.getErrAsZapFields(event), 1)
	assert.Equal(t, "MyErr", entry.getErrAsZapFields(event)[0].Key)
}

func TestGetCounterAsZapFieldsWithEmptyErrs(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent().(*EventImpl)

	assert.Empty(t, entry.getCounterAsZapFields(event))
}

func TestGetCounterAsZapFieldsHappyCase(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent().(*EventImpl)
	event.SetCounter("counter", 1)

	assert.NotEmpty(t, entry.getCounterAsZapFields(event))
	assert.Len(t, entry.getCounterAsZapFields(event), 1)
	assert.Equal(t, "counter", entry.getCounterAsZapFields(event)[0].Key)
}

func TestGetTimerAsZapFieldsWithEmptyErrs(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent().(*EventImpl)

	assert.Empty(t, entry.getTimerAsZapFields(event))
}

func TestGetTimerAsZapFieldsHappyCase(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent().(*EventImpl)
	event.SetStartTimeMS(1)
	event.StartTimer("timer")
	event.EndTimer("timer")

	assert.NotEmpty(t, entry.getTimerAsZapFields(event))
	assert.Len(t, entry.getTimerAsZapFields(event), 2)
	assert.Equal(t, "timer.elapsed_ms", entry.getTimerAsZapFields(event)[0].Key)
	assert.Equal(t, "timer.count", entry.getTimerAsZapFields(event)[1].Key)
}

func TestApplyRkFormatHappyCase(t *testing.T) {
	entry := &eventEntryImpl{}
	event := getEvent().(*EventImpl)
	event.StartTimer("timer")
	event.EndTimer("timer")
	event.AddKv("key", "value")
	event.AddErr(MyErr{})

	builder := &bytes.Buffer{}
	entry.applyRkFormat(builder, event)
	res := builder.String()

	assert.True(t, strings.Contains(res, "end_time="))
	assert.True(t, strings.Contains(res, "start_time="))
	assert.True(t, strings.Contains(res, "time="))
	assert.True(t, strings.Contains(res, "hostname="))
	assert.True(t, strings.Contains(res, "app="))
	assert.True(t, strings.Contains(res, "operation="))
	assert.True(t, strings.Contains(res, "event_status="))
}

func getEvent() Event {
	return NewEventImpl(
		&RealTimeSource{},
		"test-app",
		"test-host",
		make(map[string]string),
		zap.NewNop(),
		RK,
		false,
		make([]eventEntryListener, 0),
		false)
}

type MyErr struct{}

func (err MyErr) Error() string {
	return ""
}
