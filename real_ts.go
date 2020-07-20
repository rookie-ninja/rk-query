// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

import "time"

// Implementation of TimeSource with UnixNano() method.
type RealTimeSource struct{}

// A little bit expensive calling UnixNano() every time.
// Needs an optimization.
func (source *RealTimeSource) CurrentTimeMS() int64 {
	return time.Now().UnixNano()/(int64(time.Millisecond))
}