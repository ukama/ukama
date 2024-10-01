/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package util

import (
	"fmt"
	"strings"
)

type Config struct {
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
