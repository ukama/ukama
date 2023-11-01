/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// ProcessVersionArgument checks whether the version argument
// is present and if yes prints version and exits with 0 code
func ProcessVersionArgument(serviceName string, args []string, version string) {
	if len(os.Args) == 2 && (strings.EqualFold(os.Args[1], "version") || strings.EqualFold(strings.TrimLeft(os.Args[1], "-"), "version")) {
		fmt.Println(serviceName + " Version: " + version)
		os.Exit(0)
	}

	logrus.Infof("Starting " + serviceName + " Version: " + version)
}
