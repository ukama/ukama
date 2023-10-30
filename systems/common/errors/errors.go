/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

// This package is used to replace usage of github.com/pkg/errors package that used all over the code base
// with more conventional wrapping
package errors

import (
	"fmt"
)

func Wrap(err error, message string) error {
	return fmt.Errorf("%s %w", message, err)
}
