/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type config struct {
}

type param struct {
	value  string
	values []string
}

func (p *param) String() string {
	return string(p.value)
}

func (p *param) Set(s string) error {
	for _, v := range p.values {
		if strings.ToLower(s) == strings.ToLower(v) {
			p.value = strings.ToLower(v)
			return nil
		}
	}
	return fmt.Errorf("must match one of the following: %q", p.values)
}

func run(dir string, out io.Writer, cfg *config) error {
	return walkAndParse(dir, out)
}

func main() {
	root := flag.String("root", ".", "Root directory to start from")
	flag.Parse()

	err := run(*root, os.Stdout, &config{})
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error walking through the directory: %v\n", err)

		os.Exit(1)
	}
}
