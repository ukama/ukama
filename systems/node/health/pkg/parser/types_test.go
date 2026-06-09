/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package parser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHealthPayloadUnixReportedAt(t *testing.T) {
	payload, err := ParseHealthPayload(json.RawMessage(`{
		"schemaVersion": "1",
		"nodeType": "HomeNode",
		"reportedAt": "1779534357"
	}`))
	require.NoError(t, err)
	assert.Equal(t, int64(1779534357), payload.ReportedAt)

	payload, err = ParseHealthPayload(json.RawMessage(`{
		"schemaVersion": "1",
		"nodeType": "HomeNode",
		"reportedAt": 1779534357
	}`))
	require.NoError(t, err)
	assert.Equal(t, int64(1779534357), payload.ReportedAt)
}
