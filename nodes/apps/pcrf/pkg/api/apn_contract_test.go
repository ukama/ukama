/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package api

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestAPNLengthBoundaries(t *testing.T) {
	longest := strings.Repeat("a", 63) + "." + strings.Repeat("b", 35)
	for _, text := range []string{"a", "a-b.c9", strings.Repeat("a", 63), longest} {
		value, err := NormalizeAPN(strings.ToUpper(text))
		require.NoError(t, err)
		require.Equal(t, text, value)
	}
	for _, text := range []string{longest + "b", strings.Repeat("a", 64),
		"a..b", "a.\x00b", "é", "a_b", "a\tb", "a\nb", "-a", "a-", ".a", "a."} {
		_, err := NormalizeAPN(text)
		require.Error(t, err, text)
	}
}
