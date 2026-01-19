/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ukama_test

import (
	"database/sql/driver"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
)

func TestNodeType(t *testing.T) {
	t.Run("Constants", func(tt *testing.T) {
		assert.Equal(t, "hnode", ukama.NODE_TYPE_HOMENODE)
		assert.Equal(t, "tnode", ukama.NODE_TYPE_TOWERNODE)
		assert.Equal(t, "anode", ukama.NODE_TYPE_AMPNODE)
		assert.Equal(t, "undef", ukama.NODE_TYPE_UNDEFINED)
	})

	t.Run("StringReturnsCorrectValue", func(tt *testing.T) {
		nodeType := ukama.NodeType("hnode")
		assert.Equal(t, "hnode", nodeType.String())
	})

	t.Run("StringWithDifferentTypes", func(tt *testing.T) {
		assert.Equal(t, "hnode", ukama.NodeType(ukama.NODE_TYPE_HOMENODE).String())
		assert.Equal(t, "tnode", ukama.NodeType(ukama.NODE_TYPE_TOWERNODE).String())
		assert.Equal(t, "anode", ukama.NodeType(ukama.NODE_TYPE_AMPNODE).String())
		assert.Equal(t, "undef", ukama.NodeType(ukama.NODE_TYPE_UNDEFINED).String())
	})

	t.Run("StringWithEmptyValue", func(tt *testing.T) {
		nodeType := ukama.NodeType("")
		assert.Equal(t, "", nodeType.String())
	})

	t.Run("ValueReturnsString", func(tt *testing.T) {
		nodeType := ukama.NodeType("hnode")
		value, err := nodeType.Value()

		assert.NoError(t, err)
		assert.NotNil(t, value)
		assert.Equal(t, "hnode", value.(string))
	})

	t.Run("ValueWithDifferentTypes", func(tt *testing.T) {
		testCases := []struct {
			name     string
			nodeType ukama.NodeType
			expected string
		}{
			{"HomeNode", ukama.NodeType(ukama.NODE_TYPE_HOMENODE), "hnode"},
			{"TowerNode", ukama.NodeType(ukama.NODE_TYPE_TOWERNODE), "tnode"},
			{"AmpNode", ukama.NodeType(ukama.NODE_TYPE_AMPNODE), "anode"},
			{"Undefined", ukama.NodeType(ukama.NODE_TYPE_UNDEFINED), "undef"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(tt *testing.T) {
				value, err := tc.nodeType.Value()
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, value.(string))
			})
		}
	})

	t.Run("ValueImplementsDriverValue", func(tt *testing.T) {
		nodeType := ukama.NodeType("hnode")
		var _ driver.Value = nodeType
		value, err := nodeType.Value()
		assert.NoError(t, err)
		assert.IsType(t, "", value)
	})

	t.Run("ScanWithValidString", func(tt *testing.T) {
		var nodeType ukama.NodeType
		err := nodeType.Scan("hnode")

		assert.NoError(t, err)
		assert.Equal(t, ukama.NodeType("hnode"), nodeType)
	})

	t.Run("ScanWithDifferentTypes", func(tt *testing.T) {
		testCases := []struct {
			name     string
			input    interface{}
			expected ukama.NodeType
		}{
			{"HomeNode", "hnode", ukama.NodeType(ukama.NODE_TYPE_HOMENODE)},
			{"TowerNode", "tnode", ukama.NodeType(ukama.NODE_TYPE_TOWERNODE)},
			{"AmpNode", "anode", ukama.NodeType(ukama.NODE_TYPE_AMPNODE)},
			{"Undefined", "undef", ukama.NodeType(ukama.NODE_TYPE_UNDEFINED)},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(tt *testing.T) {
				var nodeType ukama.NodeType
				err := nodeType.Scan(tc.input)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, nodeType)
			})
		}
	})

	t.Run("ScanWithEmptyString", func(tt *testing.T) {
		var nodeType ukama.NodeType
		err := nodeType.Scan("")

		assert.NoError(t, err)
		assert.Equal(t, ukama.NodeType(""), nodeType)
	})

	t.Run("ScanWithNil", func(tt *testing.T) {
		var nodeType ukama.NodeType

		// The current implementation panics on nil
		assert.Panics(t, func() {
			_ = nodeType.Scan(nil)
		})
	})

	t.Run("ScanWithInvalidType", func(tt *testing.T) {
		var nodeType ukama.NodeType

		// The current implementation panics on non-string types
		assert.Panics(t, func() {
			_ = nodeType.Scan(123)
		})
	})

	t.Run("ScanWithInt64", func(tt *testing.T) {
		var nodeType ukama.NodeType

		// The current implementation panics on non-string types
		assert.Panics(t, func() {
			_ = nodeType.Scan(int64(123))
		})
	})

	t.Run("StringToNodeType", func(tt *testing.T) {
		str := "hnode"
		nodeType := ukama.NodeType(str)
		assert.Equal(t, "hnode", string(nodeType))
	})

	t.Run("NodeTypeToString", func(tt *testing.T) {
		nodeType := ukama.NodeType("tnode")
		str := string(nodeType)
		assert.Equal(t, "tnode", str)
	})

	t.Run("NodeTypeComparison", func(tt *testing.T) {
		nodeType1 := ukama.NodeType("hnode")
		nodeType2 := ukama.NodeType("hnode")
		nodeType3 := ukama.NodeType("tnode")

		assert.Equal(t, nodeType1, nodeType2)
		assert.NotEqual(t, nodeType1, nodeType3)
	})

	t.Run("ScanAndValueRoundTrip", func(tt *testing.T) {
		original := ukama.NodeType("hnode")

		// Convert to driver.Value
		value, err := original.Value()
		assert.NoError(t, err)

		// Scan back from driver.Value
		var scanned ukama.NodeType
		err = scanned.Scan(value)
		assert.NoError(t, err)

		// Should be equal
		assert.Equal(t, original, scanned)
	})

	t.Run("StringAndValueConsistency", func(tt *testing.T) {
		nodeType := ukama.NodeType("anode")

		stringValue := nodeType.String()
		driverValue, err := nodeType.Value()
		assert.NoError(t, err)

		assert.Equal(t, stringValue, driverValue.(string))
	})
}
