/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfigurationStateDefault(t *testing.T) {
	config := NewConfig("controller")
	require.Equal(t, "state:9090", config.StateHost)
	require.Equal(t, "controller", config.DB.DbName)
}
