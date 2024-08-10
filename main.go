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

var oFile = os.Stdout

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
	src := flag.String("src", ".", "Source directory to start from")
	ofilename := flag.String("f", "", "The name of the file to write to (or stdout)")
	flag.Parse()

	var err error
	if *ofilename != "" {
		oFile, err = os.Create(*ofilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer oFile.Close()
	}

	err = run(*src, oFile, &config{})
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error walking through the directory: %v\n", err)

		os.Exit(1)
	}
}
