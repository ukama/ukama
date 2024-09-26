/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package scan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/ukama/msgcli/util"
)

type Config struct {
	OutputFormat string
}

type Param struct {
	Value  string
	Values []string
}

func (p *Param) String() string {
	return string(p.Value)
}

func (p *Param) Set(s string) error {
	for _, v := range p.Values {
		if strings.EqualFold(s, v) {
			p.Value = strings.ToLower(v)
			return nil
		}
	}
	return fmt.Errorf("must match one of the following: %q", p.Values)
}

func (p *Param) Type() string {
	return "String"
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

func Run(dir string, out io.Writer, cfg *Config) error {
	data := &util.ResultSet{}

	err := WalkAndParse(dir, data)
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
