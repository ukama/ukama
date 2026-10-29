/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package engine_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ukama/ukama/systems/analytics/ingest/pkg/engine"
)

func decode(t *testing.T, s string) interface{} {
	t.Helper()

	var v interface{}
	require.NoError(t, json.Unmarshal([]byte(s), &v))

	return v
}

func TestGetPath(t *testing.T) {
	item := decode(t, `{"metric":{"network":"n1","node_id":"uk-1"},"value":[1793169036.5,"42"],"grid":[[1,2],[3,4]]}`)

	tests := []struct {
		name  string
		path  string
		want  interface{}
		found bool
	}{
		{"root", "$", item, true},
		{"nested key", "$.metric.network", "n1", true},
		{"normalized key", "$.metric.nodeId", "uk-1", true},
		{"whole array", "$.value", []interface{}{1793169036.5, "42"}, true},
		{"array index", "$.value[1]", "42", true},
		{"first index", "$.value[0]", 1793169036.5, true},
		{"nested index", "$.grid[1][0]", float64(3), true},
		{"index out of range", "$.value[2]", nil, false},
		{"index on non-array", "$.metric[0]", nil, false},
		{"missing key", "$.metric.site", nil, false},
		{"negative index", "$.value[-1]", nil, false},
		{"non-numeric index", "$.value[x]", nil, false},
		{"unclosed bracket", "$.value[1", nil, false},
		{"junk after index", "$.value[1]x", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := engine.GetPath(item, tt.path)
			assert.Equal(t, tt.found, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHashIgnoresSampleTimestamp(t *testing.T) {
	fieldMap := map[string]string{"network_id": "$.metric.network", "value": "$.value[1]"}

	first := engine.MapItem(decode(t, `{"metric":{"network":"n1"},"value":[1793169036,"42"]}`), fieldMap, nil)
	later := engine.MapItem(decode(t, `{"metric":{"network":"n1"},"value":[1793169096,"42"]}`), fieldMap, nil)
	changed := engine.MapItem(decode(t, `{"metric":{"network":"n1"},"value":[1793169156,"43"]}`), fieldMap, nil)

	assert.Equal(t, engine.HashFields(first), engine.HashFields(later))
	assert.NotEqual(t, engine.HashFields(first), engine.HashFields(changed))
}
