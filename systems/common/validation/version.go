/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package validation

import (
	"errors"

	"github.com/coreos/go-semver/semver"
)

func ParseVersion(version string) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, errors.New("Invalid version format. Refer to https://semver.org/ for more information")
	}

	return v, err
}

func CompareVersions(version1 string, version2 string) (int, error) {
	v1, err := ParseVersion(version1)
	if err != nil {
		return 0, err
	}
	v2, err := ParseVersion(version2)
	if err != nil {
		return 0, err
	}
	c := v1.Compare(*v2)
	return c, nil // 1 if version1 is greater than version2, 0 if they are equal, -1 if version1 is less than version2
}

