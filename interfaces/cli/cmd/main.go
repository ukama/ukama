/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"github.com/spf13/cobra"
	"github.com/ukama/ukama/interfaces/cli/pkg/cmd"
)

func main() {
	cobra.CheckErr(cmd.RootCommand().Execute())
}
