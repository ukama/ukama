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

type ReportType uint8

const (
	ReportTypeUnknown ReportType = iota
	ReportTypeInvoice
	ReportTypeConsumption
)

func (r *ReportType) Scan(value interface{}) error {
	*r = ReportType(uint8(value.(int64)))

	return nil
}

func (r ReportType) Value() (driver.Value, error) {
	return int64(r), nil
}

func (r ReportType) String() string {
	t := map[ReportType]string{0: "unknown", 1: "invoice", 2: "consumption"}

	v, ok := t[r]
	if !ok {
		return t[0]
	}

	return v
}

func ParseReportType(value string) ReportType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return ReportType(i)
	}

	t := map[string]ReportType{"unknown": 0, "invoice": 1, "consumption": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return ReportType(0)
	}

	return ReportType(v)
}
