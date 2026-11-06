/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"os"
	"testing"
)

func TestConfiguredAPN(t *testing.T) {
	t.Setenv("UKAMA_APN", "")
	require.NoError(t, os.Unsetenv("UKAMA_APN"))
	require.Equal(t, "internet", NewConfig("pcrf").APN)
	t.Setenv("UKAMA_APN", "Corp.Ukama")
	got, err := api.NormalizeAPN(NewConfig("pcrf").APN)
	require.NoError(t, err)
	require.Equal(t, "corp.ukama", got)
	for _, value := range []string{"", " internet", "a..b"} {
		t.Setenv("UKAMA_APN", value)
		_, err = api.NormalizeAPN(NewConfig("pcrf").APN)
		require.Error(t, err, value)
	}
}
