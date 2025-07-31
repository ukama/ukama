/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = make(JSONMap)
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONMap: %v", value)
	}
	return json.Unmarshal(data, m)
}
