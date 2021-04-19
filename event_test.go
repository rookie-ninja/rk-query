package rkquery

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEventFromContext_WithNilContext(t *testing.T) {
	event := GetEventFromContext(nil)

	assert.NotNil(t, event)
}

func TestGetEventFromContext_WithEmptyEvent(t *testing.T) {
	event := GetEventFromContext(context.Background())

	assert.NotNil(t, event)
}

func TestGetEventFromContext_HappyCase(t *testing.T) {
	event := &eventNoop{}
	ctx := context.WithValue(context.Background(), EventKey, event)

	assert.Equal(t, event, GetEventFromContext(ctx))
}