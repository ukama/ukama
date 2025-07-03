/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

type CdrType uint8

const (
	CdrTypeUnknown CdrType = iota
	CdrTypeData
	CdrTypeVoice
	CdrTypeSms
	CdrTypeVoicemail
	CdrTypeDataUnlimited
	CdrTypeDataCaptivePortal
)

func (s *CdrType) Scan(value interface{}) error {
	*s = CdrType(uint8(value.(int64)))

	return nil
}

func (s CdrType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s CdrType) String() string {
	t := map[CdrType]string{0: "unknown", 1: "data", 2: "voice", 3: "sms", 4: "voicemail",
		5: "data_unlimited", 6: "data_captive_portal"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseCdrType(value string) CdrType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return CdrType(i)
	}

	t := map[string]CdrType{"unknown": 0, "data": 1, "voice": 2, "sms": 3, "voicemail": 4,
		"data_unlimited": 5, "data_captive_portal": 6}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return CdrType(0)
	}

	return CdrType(v)
}
