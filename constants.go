// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkquery

const (
	scopeDelimiter = "------------------------------------------------------------------------"
	eoe            = "EOE"
	unknown        = "unknown"
	openMarker     = "-open-"
	// ************* Time *************
	startTimeKey = "startTime"
	endTimeKey   = "endTime"
	elapsedKey   = "elapsedNano"
	timezoneKey  = "timezone"
	// ************* App *************
	serviceKey        = "service"
	serviceNameKey    = "serviceName"
	serviceVersionKey = "serviceVersion"
	entryNameKey      = "entryName"
	entryKindKey      = "entryKind"
	// ************* Env *************
	envKey      = "env"
	hostnameKey = "hostname"
	localIpKey  = "localIP"
	domainKey   = "domain"
	goosKey     = "os"
	goArchKey   = "arch"
	// ************* Ids *************
	idsKey       = "ids"
	eventIdKey   = "eventId"
	traceIdKey   = "traceId"
	requestIdKey = "requestId"
	// ************* Payloads *************
	payloadsKey = "payloads"
	// ************* Counters *************
	countersKey = "counters"
	// ************* Pairs *************
	pairsKey       = "pairs"
	resCodeKey     = "resCode"
	operationKey   = "operation"
	remoteAddrKey  = "remoteAddr"
	eventStatusKey = "eventStatus"
	timingKey      = "timing"
	errKey         = "error"
)
