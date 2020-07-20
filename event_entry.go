package rk_query

import "go.uber.org/zap"

// Immutable event entries.
type eventEntry interface {
	FormatAsRk(Event) string

	FormatAsJson(Event) *zap.Logger
}