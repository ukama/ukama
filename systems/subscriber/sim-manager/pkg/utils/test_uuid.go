/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

// This is a slightly  modified version of
// https://github.com/google/uuid/tree/v1.0.0

// The 3-BSD clause allows for modification, private use
// commercial use and distribution as indicated at
// https://github.com/google/uuid/blob/v1.0.0/LICENSE

// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"errors"
	"fmt"
	"strings"
)

const testUUUIDPrefix = "testuuid"

// xvalues returns the value of a byte as a hexadecimal digit or 255.
var xvalues = [256]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
}

// ParseTestUUID return whether or not the provided string matches a test UUID.  Both the UUID form of
// testuuid-xxxx-xxxx-xxxx-xxxxxxxxxxxx and
// urn:uuid:testuuid-xxxx-xxxx-xxxx-xxxxxxxxxxxx are supported.
func ParseTestUUID(s string) error {
	if len(s) != 36 {
		if len(s) != 36+9 {
			return fmt.Errorf("invalid test UUID length: %d", len(s))
		}
		if strings.ToLower(s[:9]) != "urn:uuid:" {
			return fmt.Errorf("invalid urn prefix: %q", s[:9])
		}
		s = s[9:]
	}

	if !strings.HasPrefix(s, testUUUIDPrefix) {
		return fmt.Errorf("invalid test uuid prefix: %q", s[:9])
	}

	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return errors.New("invalid test UUID format")
	}
	for _, x := range [12]int{
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34} {
		_, ok := xtob(s[x], s[x+1])
		if !ok {
			return errors.New("invalid test UUID format")
		}
	}
	return nil
}

func GetIccidFromTestUUID(s string) (string, error) {
	if err := ParseTestUUID(s); err != nil {
		return "", err
	}

	iccid := strings.ReplaceAll(strings.TrimPrefix(s, testUUUIDPrefix), "-", "")

	return iccid, nil
}

// xtob converts hex characters x1 and x2 into a byte.
func xtob(x1, x2 byte) (byte, bool) {
	b1 := xvalues[x1]
	b2 := xvalues[x2]
	return (b1 << 4) | b2, b1 != 255 && b2 != 255
}
