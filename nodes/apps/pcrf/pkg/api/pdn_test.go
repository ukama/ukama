/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeAPN(t *testing.T) {
	for _, value := range []string{"internet", "INTERNET", "ukama.example"} {
		got, err := NormalizeAPN(value)
		require.NoError(t, err)
		require.Equal(t, strings.ToLower(value), got)
	}
	for _, value := range []string{"", " internet", "internet ", "a..b",
		"-a", "a-", "a/b", "a\x00b", "ä", "a_b", ".a", "a.",
		strings.Repeat("a", 64), strings.Repeat("a.", 50) + "b"} {
		_, err := NormalizeAPN(value)
		require.Error(t, err, value)
	}
}
