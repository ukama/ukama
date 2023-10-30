/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"errors"
	"time"
)

func ValidateDOB(dob string) (string, error) {
	t, err := time.Parse(time.RFC3339, dob)
	if err != nil {
		return "", errors.New("invalid date format, must be RFC3339 standard")
	}
	return t.Format(time.RFC3339), nil
}
