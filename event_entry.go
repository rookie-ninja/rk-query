// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import "go.uber.org/zap"

// Currently, we support RK format and Zap logger json format
// More logger option would be applied
type eventEntry interface {
	FormatAsRk(Event) string

	FormatAsRkMin(Event) string

	FormatAsZap(Event) *zap.Logger

	FormatAsZapMin(Event) *zap.Logger
}