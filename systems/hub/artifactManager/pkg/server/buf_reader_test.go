/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var str = []byte("the quick brown fox jumps over the lazy dog")

func Test_Read(t *testing.T) {
	reader := NewBufReader(bytes.NewReader(str))
	b, err := io.ReadAll(reader)

	assert.NoError(t, err)
	assert.Equal(t, str, b)
}

func Test(t *testing.T) {
	tests := []struct {
		name     string
		preReads func(io.Reader)
	}{
		{
			name: "PreReadByte",
			preReads: func(r io.Reader) {
				_, err := r.Read(make([]byte, 1))
				assert.NoError(t, err)
			},
		},
		{
			name: "PreReadAll",
			preReads: func(r io.Reader) {
				_, err := r.Read(make([]byte, len(str)))
				assert.NoError(t, err)
			},
		},
		{
			name: "PreReadNothing",
			preReads: func(r io.Reader) {
				_, err := r.Read(make([]byte, 0))
				assert.NoError(t, err)
			},
		},
		{
			name: "PreReadMoreThenStream",
			preReads: func(r io.Reader) {
				_, err := r.Read(make([]byte, 1024))
				assert.NoError(t, err)
			},
		},
		{
			name: "PreReadTwice",
			preReads: func(r io.Reader) {
				_, err := r.Read(make([]byte, 1))
				assert.NoError(t, err)

				_, err = r.Read(make([]byte, 2))
				assert.NoError(t, err)
			},
		},
		{
			name: "PreReadAllByOne",
			preReads: func(r io.Reader) {
				for i := 0; i < len(str); i++ {
					_, err := r.Read(make([]byte, 1))
					assert.NoError(t, err)
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := NewBufReader(bytes.NewReader(str))

			test.preReads(reader)
			reader.Reset()

			b, err := io.ReadAll(reader)

			assert.NoError(t, err)
			assert.Equal(t, str, b)
		})
	}
}

func Test_ResetTwice(t *testing.T) {
	t.Run("RestOnce", func(t *testing.T) {
		reader := NewBufReader(bytes.NewReader(str))
		reader.Reset()
		assert.Panics(t, func() {
			reader.Reset()
		})
	})
}
