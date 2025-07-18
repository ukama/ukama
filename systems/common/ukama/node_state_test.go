/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
)

func TestNodeState(t *testing.T) {
	t.Run("NodeStateValidString", func(tt *testing.T) {
		configuredState := ukama.ParseNodeState("configureD")

		assert.NotNil(t, configuredState)
		assert.Equal(t, configuredState.String(), ukama.NodeStateConfigured.String())
		assert.Equal(t, uint8(configuredState), uint8(1))
		assert.True(t, ukama.IsValidNodeState("configureD"))
	})

	t.Run("NodeStateValidNumber", func(tt *testing.T) {
		operationalState := ukama.ParseNodeState("2")

		assert.NotNil(t, operationalState)
		assert.Equal(t, uint8(operationalState), uint8(2))
		assert.Equal(t, operationalState.String(), ukama.NodeStateOperational.String())
		assert.True(t, ukama.IsValidNodeState("operational"))
	})

	t.Run("NodeStateNonValidString", func(tt *testing.T) {
		unknownState := ukama.ParseNodeState("failure")

		assert.NotNil(t, unknownState)
		assert.Equal(t, unknownState.String(), ukama.NodeStateUnknown.String())
		assert.Equal(t, uint8(unknownState), uint8(0))
		assert.False(t, ukama.IsValidNodeState("failure"))
	})

	t.Run("NodeStateNonValidNumber", func(tt *testing.T) {
		unknownState := ukama.NodeState(uint8(10))

		assert.NotNil(t, unknownState)
		assert.Equal(t, unknownState.String(), ukama.NodeStateUnknown.String())
		assert.Equal(t, uint8(unknownState), uint8(10))
	})

	t.Run("NodeStateNonValidStrNumber", func(tt *testing.T) {
		unknownState := ukama.ParseNodeState("10")

		assert.NotNil(t, unknownState)
		assert.Equal(t, unknownState.String(), ukama.NodeStateUnknown.String())
		assert.Equal(t, uint8(unknownState), uint8(10))
	})
}
