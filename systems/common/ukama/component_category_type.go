// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2023-present, Ukama Inc.

package ukama

import "strconv"

type ComponentCategory uint8

// TODO: missing sentinel value for all non valid component categories of the universe.
const (
	// TODO: golang fmt standard advises not using all caps as constants names.
	// These should be camel or title cases instead.
	ALL ComponentCategory = iota
	ACCESS
	BACKHAUL
	POWER
	SWITCH
	SPECTRUM
)

func (c *ComponentCategory) Scan(value interface{}) error {
	*c = ComponentCategory(uint8(value.(int64)))
	return nil
}

func (c ComponentCategory) Value() (uint8, error) {
	return uint8(c), nil
}

func (c ComponentCategory) String() string {
	t := map[ComponentCategory]string{
		0: "all",
		1: "access",
		2: "backhaul",
		3: "power",
		4: "switch",
		5: "spectrum",
	}

	v, ok := t[c]
	if !ok {
		return t[0]
	}

	return v
}

// TODO: this should be renamed to ParseComponentCategory instead + update
// all dependents services using it.
// TODO: make this also case insensitive.
func ParseType(value string) ComponentCategory {
	i, err := strconv.Atoi(value)
	if err == nil {
		return ComponentCategory(i)
	}

	t := map[string]ComponentCategory{
		"all":      0,
		"access":   1,
		"backhaul": 2,
		"power":    3,
		"switch":   4,
		"spectrum": 5,
	}

	v, ok := t[value]
	if !ok {
		return ComponentCategory(0)
	}

	return ComponentCategory(v)
}
