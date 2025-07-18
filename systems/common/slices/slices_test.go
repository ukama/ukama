/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	ID int
}

func Test_FindStruct(t *testing.T) {
	s := []Foo{
		Foo{ID: 1},
		Foo{ID: 2},
		Foo{ID: 3},
	}

	t.Run("ItemFound", func(tt *testing.T) {
		f := Find(s, func(i *Foo) bool {
			return i.ID == 2
		})

		assert.NotNil(tt, f)
		assert.Equal(tt, 2, f.ID)
	})

	t.Run("ItemNotFound", func(tt *testing.T) {
		f := Find(s, func(i *Foo) bool {
			return i.ID == 5
		})

		assert.Nil(tt, f)
	})
}

func Test_FindPointer(t *testing.T) {
	s := []*Foo{
		&Foo{ID: 1},
		&Foo{ID: 2},
		&Foo{ID: 3},
	}

	t.Run("ItemFound", func(tt *testing.T) {
		f := FindPointer(s, func(i *Foo) bool {
			return i.ID == 2
		})

		assert.NotNil(tt, f)
		assert.Equal(t, 2, f.ID)
	})

	t.Run("ItemNotFound", func(tt *testing.T) {
		f := FindPointer(s, func(i *Foo) bool {
			return i.ID == 10
		})

		assert.Nil(tt, f)
	})
}
