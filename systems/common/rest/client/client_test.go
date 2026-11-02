/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	crest "github.com/ukama/ukama/systems/common/rest"
)

type customError struct {
	Message string `json:"message"`
}

func (e customError) Error() string {
	return e.Message
}

func TestNewResty(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		r := NewResty()

		assert.False(t, r.C.Debug)
		assert.Equal(t, time.Duration(0), r.C.GetClient().Timeout)
		assert.Equal(t, reflect.TypeOf(crest.ErrorResponse{}), r.C.Error)
		assert.Empty(t, r.C.Header.Get("Authorization"))
		assert.Empty(t, r.C.Header.Get("Content-Type"))
		assert.Empty(t, r.C.Header.Get("Accept"))
	})

	t.Run("WithDebug", func(t *testing.T) {
		assert.True(t, NewResty(WithDebug(true)).C.Debug)
		assert.False(t, NewResty(WithDebug(false)).C.Debug)
	})

	t.Run("WithBearer", func(t *testing.T) {
		r := NewResty(WithBearer("secret"))

		assert.Equal(t, "Bearer secret", r.C.Header.Get("Authorization"))
	})

	t.Run("WithError", func(t *testing.T) {
		r := NewResty(WithError(customError{}))

		assert.Equal(t, reflect.TypeOf(customError{}), r.C.Error)
	})

	t.Run("WithContentTypeJSON", func(t *testing.T) {
		r := NewResty(WithContentTypeJSON())

		assert.Equal(t, "application/json", r.C.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.C.Header.Get("Accept"))
	})

	t.Run("WithTimeout", func(t *testing.T) {
		r := NewResty(WithTimeout(3 * time.Second))

		assert.Equal(t, 3*time.Second, r.C.GetClient().Timeout)
	})

	t.Run("CombinedOptions", func(t *testing.T) {
		r := NewResty(WithDebug(true), WithBearer("secret"), WithTimeout(time.Second))

		assert.True(t, r.C.Debug)
		assert.Equal(t, "Bearer secret", r.C.Header.Get("Authorization"))
		assert.Equal(t, time.Second, r.C.GetClient().Timeout)
		assert.Equal(t, reflect.TypeOf(crest.ErrorResponse{}), r.C.Error)
	})
}
