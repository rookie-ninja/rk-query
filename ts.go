// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

// Interface for getting current time which could be implemented by users
// This is mainly for testing purpose.
type TimeSource interface {
	// get the current time in milliseconds.
	CurrentTimeMS() int64
}