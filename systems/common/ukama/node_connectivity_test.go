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

func TestNodeConnectivity(t *testing.T) {
	t.Run("NodeConnectivityValidString", func(tt *testing.T) {
		onlineConnectivity := ukama.ParseNodeConnectivity("onlinE")

		assert.NotNil(t, onlineConnectivity)
		assert.Equal(t, onlineConnectivity.String(), "online")
		assert.Equal(t, uint8(onlineConnectivity), uint8(1))
	})

	t.Run("NodeConnectivityValidNumber", func(tt *testing.T) {
		offlineConnectivity := ukama.ParseNodeConnectivity("2")

		assert.NotNil(t, offlineConnectivity)
		assert.Equal(t, uint8(offlineConnectivity), uint8(2))
		assert.Equal(t, offlineConnectivity.String(), "offline")
	})

	t.Run("NodeConnectivityNonValidString", func(tt *testing.T) {
		unsupportedConnectivity := ukama.ParseNodeConnectivity("failure")

		assert.NotNil(t, unsupportedConnectivity)
		assert.Equal(t, unsupportedConnectivity.String(), "unknown")
		assert.Equal(t, uint8(unsupportedConnectivity), uint8(0))
	})

	t.Run("NodeConnectivityNonValidNumber", func(tt *testing.T) {
		unknownConnecitivty := ukama.NodeConnectivity(uint8(10))

		assert.NotNil(t, unknownConnecitivty)
		assert.Equal(t, unknownConnecitivty.String(), "unknown")
		assert.Equal(t, uint8(unknownConnecitivty), uint8(10))
	})
}
