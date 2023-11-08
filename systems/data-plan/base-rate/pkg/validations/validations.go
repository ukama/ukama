/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package validations

func IsEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func IsValidUploadReqArgs(fileUrl, effectiveAt, simType string) bool {
	return !IsEmpty(fileUrl, effectiveAt, simType)
}
