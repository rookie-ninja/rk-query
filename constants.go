// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

const (
	EventScopeDelimiter       string = "------------------------------------------------------------------------"
	EventEOE                  string = "EOE"
	EventUnknownApplication   string = "Unknown Application"
	EventUnknownHostName      string = "Unknown Hostname"
	EventTruncatedString      string = "TRUNCATED"
	EventCommaTruncatedString string = ",TRUNCATED"
	EventMaxHistoryLength     int    = 1024
	EventOperationKey         string = "operation"
	EventStatusKey            string = "status"
	EventOpenMarker           string = "-open-"
)