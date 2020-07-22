// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import "go.uber.org/zap"

// Immutable event entries.
type eventEntry interface {
	FormatAsRk(Event) string

	FormatAsJson(Event) *zap.Logger
}