// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_query

const (
	ScopeDelimiter   string = "------------------------------------------------------------------------"
	EOE              string = "EOE"
	Unknown          string = "Unknown"
	Truncated        string = "TRUNCATED"
	CommaTruncated   string = ",TRUNCATED"
	MaxHistoryLength int    = 1024
	OpenMarker       string = "-open-"
	appNameKey              = "app_name"
	hostnameKey             = "hostname"
	operationKey            = "operation"
	remoteAddrKey           = "remote_addr"
	eventStatusKey          = "event_status"
	historyKey              = "history"
	startTimeKey            = "start_time"
	endTimeKey              = "end_time"
	timeKey                 = "time"
	timingKey               = "timing"
	counterKey              = "counter"
	pairKey                 = "pair"
	errKey                  = "error"
	fieldKey                = "field"
)
