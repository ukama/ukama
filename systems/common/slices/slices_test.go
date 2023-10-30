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

	f := Find(s, func(i *Foo) bool {
		return i.ID == 2
	})

	assert.Equal(t, 2, f.ID)
}

func Test_FindPointer(t *testing.T) {
	s := []*Foo{
		&Foo{ID: 1},
		&Foo{ID: 2},
		&Foo{ID: 3},
	}

	f := FindPointer(s, func(i *Foo) bool {
		return i.ID == 2
	})

	assert.Equal(t, 2, f.ID)
}
