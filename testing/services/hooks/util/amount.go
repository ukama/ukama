/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package util

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	errFormatMsg = "%w: %s"
)

var (
	ErrParseExp = errors.New("failed to parse regex for amount")
	ErrMatch    = errors.New("failed to match regex for amount")
	ErrInt      = errors.New("failed to parse int for amount")
	ErrFloat    = errors.New("failed to parse float for amount")
)

func ToAmountCents(s string) (int64, error) {
	ok, err := regexp.MatchString("^([0]|([1-9][0-9]{0,17}))([.][0-9]{0,3}[1-9])?$", s)
	if err != nil {
		return 0, fmt.Errorf(errFormatMsg, ErrParseExp, err)
	}

	if !ok {
		return 0, fmt.Errorf(errFormatMsg, ErrMatch, s)
	}

	if strings.Contains(s, ".") {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return 0, fmt.Errorf(errFormatMsg, ErrFloat, err)
		}

		return int64(f * 100), nil
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(errFormatMsg, ErrInt, err)
	}

	return i * 100, nil
}

func ToAmount(amountCents int64) string {
	return strconv.FormatFloat(float64(amountCents)/100.0, 'f', 2, 64)
}
