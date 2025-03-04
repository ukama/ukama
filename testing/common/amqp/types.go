/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package amqp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OutputFormat string
}

type PushConfig struct {
	ClusterURL   string
	ClusterUsr   string
	ClusterPswd  string
	Vhost        string
	Exchange     string
	OutputFormat string
}

type EnumParam struct {
	Value  string
	Values []string
}

func (p *EnumParam) String() string {
	return string(p.Value)
}

func (p *EnumParam) Set(s string) error {
	for _, v := range p.Values {
		if strings.EqualFold(s, v) {
			p.Value = strings.ToLower(v)
			return nil
		}
	}

	return fmt.Errorf("must match one of the following: %q", p.Values)
}

func (p *EnumParam) Type() string {
	return "string"
}

type ResultSet struct {
	Routes []string
}

func Serialize(data interface{}, format string) ([]byte, error) {
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
	default:
		return nil, fmt.Errorf("specified format not supported: %v", format)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
