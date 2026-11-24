/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthClientProvider(t *testing.T) {
	t.Run("ClientIsCreatedOnceAndReused", func(t *testing.T) {
		p := NewHealthClientProvider("health:9090")

		first, err := p.GetClient()
		require.NoError(t, err)
		second, err := p.GetClient()
		require.NoError(t, err)

		assert.NotNil(t, first)
		assert.Same(t, first, second)
	})

	t.Run("InvalidHostReturnsError", func(t *testing.T) {
		p := NewHealthClientProvider("%zz")

		c, err := p.GetClient()

		assert.Error(t, err)
		assert.Nil(t, c)
	})
}
