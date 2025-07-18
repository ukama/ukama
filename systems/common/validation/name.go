/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package validation

import "regexp"

const (
	segment            = "[a-z0-9]([-_a-z0-9]*[a-z0-9])?"
	dnsLabelNameRegexp = "(" + segment + "\\.)*" + segment

	maxNameLength int = 253 // max length of DNS label
)

var sysctlRegexp = regexp.MustCompile("^" + dnsLabelNameRegexp + "$")

// IsValidDnsLabelName checks that the given string is a valid name to us in url (DNS label)
func IsValidDnsLabelName(name string) bool {
	if len(name) > maxNameLength {
		return false
	}
	return sysctlRegexp.MatchString(name)
}
