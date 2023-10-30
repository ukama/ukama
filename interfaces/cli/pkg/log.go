/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"fmt"
	"io"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type logger struct {
	stderr  io.Writer
	stdout  io.Writer
	verbose bool
}

func NewLogger(stdout io.Writer, stderr io.Writer, verbose bool) Logger {
	return &logger{
		stderr:  stderr,
		stdout:  stdout,
		verbose: verbose,
	}
}

func (l logger) Printf(format string, v ...interface{}) {
	if l.verbose {
		fmt.Fprintf(l.stdout, format, v...)
	}
}

func (l logger) Errorf(format string, v ...interface{}) {
	if l.verbose {
		fmt.Fprintf(l.stderr, format, v...)
	}
}
