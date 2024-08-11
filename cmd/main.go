/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ukama/msgcli/util"
	"gopkg.in/yaml.v3"

	"github.com/ukama/msgcli/internal/action"
)

const (
	defaultOutputFormat = "json"
)

type config struct {
	OutputFormat string
}

type param struct {
	Value  string
	Values []string
}

func (p *param) String() string {
	return string(p.Value)
}

func (p *param) Set(s string) error {
	for _, v := range p.Values {
		if strings.EqualFold(s, v) {
			p.Value = strings.ToLower(v)
			return nil
		}
	}
	return fmt.Errorf("must match one of the following: %q", p.Values)
}

var (
	oFile        = os.Stdout
	outputFormat = param{
		Values: []string{"json", "yaml", "toml"},
	}
)

func init() {
	flag.Var(&outputFormat, "oFormat",
		fmt.Sprintf("Output format. Must match one of the following: %q (default \"json\" )",
			outputFormat.Values))
}

func serialize(data interface{}, format string) (io.Writer, error) {
	var err error
	buf := &bytes.Buffer{}

	switch format {
	case "json":
		enc := json.NewEncoder(buf)
		enc.SetIndent("", "    ")
		err = enc.Encode(data)
	case "yaml":
		enc := yaml.NewEncoder(buf)
		enc.SetIndent(4)
		err = enc.Encode(data)
	case "toml":
		err = toml.NewEncoder(buf).Encode(data)
	default:
		return nil, fmt.Errorf("specified format not supported: %v", format)
	}

	return buf, err
}

func run(dir string, out io.Writer, cfg *config) error {
	data := &util.ResultSet{}

	err := action.WalkAndParse(dir, data)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	outputBuf, err := serialize(data, cfg.OutputFormat)
	if err != nil {
		return fmt.Errorf("error while writiing output: %w", err)
	}

	fmt.Fprint(out, outputBuf)

	return nil
}

func main() {
	src := flag.String("src", ".", "Source directory to start from")
	ofilename := flag.String("out", "", "The name of the file to write to (default \"Stdout\")")
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

	if outputFormat.String() == "" {
		err = outputFormat.Set(defaultOutputFormat)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	cfg := &config{
		OutputFormat: outputFormat.String(),
	}

	err = run(*src, oFile, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
