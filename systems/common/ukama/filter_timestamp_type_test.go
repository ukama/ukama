/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ukama_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
)

func TestFilterTimestampType(t *testing.T) {
	t.Run("Parse and String for all", func(tt *testing.T) {
		filter := ukama.ParseFilterTimestampType("all")

		assert.NotNil(t, filter)
		assert.Equal(t, filter.String(), ukama.FilterTimestampTypeAll.String())
		assert.Equal(t, uint8(filter), uint8(1))
		assert.Equal(t, ukama.ReturnFilterTimestampType(filter), int64(1))
	})

	t.Run("Parse and String for latest", func(tt *testing.T) {
		filter := ukama.ParseFilterTimestampType("latest")

		assert.NotNil(t, filter)
		assert.Equal(t, filter.String(), ukama.FilterTimestampTypeLatest.String())
		assert.Equal(t, uint8(filter), uint8(2))
		assert.Equal(t, ukama.ReturnFilterTimestampType(filter), int64(2))
	})

	t.Run("Parse invalid value returns unknown", func(tt *testing.T) {
		filter := ukama.ParseFilterTimestampType("invalid")

		assert.NotNil(t, filter)
		assert.Equal(t, filter.String(), ukama.FilterTimestampTypeUnknown.String())
		assert.Equal(t, uint8(filter), uint8(0))
		assert.Equal(t, ukama.ReturnFilterTimestampType(filter), int64(0))
	})

	t.Run("ReturnFilterTimestampType for out of range", func(tt *testing.T) {
		filter := ukama.FilterTimestampType(99)

		assert.NotNil(t, filter)
		assert.Equal(t, ukama.ReturnFilterTimestampType(filter), int64(0))
	})

	t.Run("Value returns int64 enum value", func(tt *testing.T) {
		v, err := ukama.FilterTimestampTypeLatest.Value()

		assert.NoError(t, err)
		assert.Equal(t, v, int64(2))
	})

	t.Run("Scan sets enum from int64 value", func(tt *testing.T) {
		var filter ukama.FilterTimestampType
		err := (&filter).Scan(int64(1))

		assert.NoError(t, err)
		assert.Equal(t, filter, ukama.FilterTimestampTypeAll)
	})
}
