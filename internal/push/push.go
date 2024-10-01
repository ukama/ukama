/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package push

import "fmt"

func Run(org, route, msg string) error {
	fmt.Print("Push called with args: ", org, route, msg)

	return nil
}
