/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package store

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	kv := NewMemKV()
	s := NewStoreWithKV(kv, 2*time.Second)

	t.Run("PutAndGet", func(t *testing.T) {
		err := s.Put("key1", "value1")
		require.NoError(t, err)
		val, err := s.Get("key1")
		require.NoError(t, err)
		assert.Equal(t, "value1", val)
	})

	t.Run("GetNotFound", func(t *testing.T) {
		_, err := s.Get("nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("PutJsonAndGetJson", func(t *testing.T) {
		type data struct {
			State  string  `json:"state"`
			Value  float64 `json:"value"`
		}
		d := data{State: "healthy", Value: 42.5}
		err := s.PutJson("json_key", d)
		require.NoError(t, err)
		bytes, err := s.GetJson("json_key")
		require.NoError(t, err)
		var decoded data
		err = json.Unmarshal(bytes, &decoded)
		require.NoError(t, err)
		assert.Equal(t, "healthy", decoded.State)
		assert.Equal(t, 42.5, decoded.Value)
	})

	t.Run("PutJsonInvalidValue", func(t *testing.T) {
		err := s.PutJson("bad", make(chan int))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "marshal")
	})

	t.Run("GetJsonNotFound", func(t *testing.T) {
		_, err := s.GetJson("missing")
		require.Error(t, err)
	})

	t.Run("GetAll", func(t *testing.T) {
		_ = s.Put("a", "1")
		_ = s.Put("b", "2")
		keys, err := s.GetAll()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(keys), 2)
		assert.Contains(t, keys, "a")
		assert.Contains(t, keys, "b")
	})

	t.Run("Delete", func(t *testing.T) {
		_ = s.Put("to_delete", "x")
		err := s.Delete("to_delete")
		require.NoError(t, err)
		_, err = s.Get("to_delete")
		require.Error(t, err)
	})
}

