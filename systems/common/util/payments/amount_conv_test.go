/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package payments

import (
	"errors"
	"testing"
)

func TestToAmountCents(t *testing.T) {
	tests := []struct {
		name        string
		amount      string
		amountCents int64
		expErr      error
	}{
		{
			name:        "positive float amount",
			amount:      "1.2",
			amountCents: 120,
			expErr:      nil,
		},

		{
			name:        "positive int amount",
			amount:      "2",
			amountCents: 200,
			expErr:      nil,
		},

		{
			name:        "positive amount with 2 decimals",
			amount:      "1.20",
			amountCents: 120,
			expErr:      nil,
		},

		{
			name:        "positive amount with 3 decimals",
			amount:      "1.200",
			amountCents: 0,
			expErr:      ErrMatch,
		},

		{
			name:        "negative amount",
			amount:      "-7",
			amountCents: 0,
			expErr:      ErrMatch,
		},

		{
			name:        "negative float amount",
			amount:      "-1.5",
			amountCents: 0,
			expErr:      ErrMatch,
		},

		{
			name:        "positive lesser than 1",
			amount:      "0.1",
			amountCents: 10,
			expErr:      nil,
		},

		{
			name:        "zero",
			amount:      "0.00",
			amountCents: 0,
			expErr:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ToAmountCents(test.amount)

			if test.expErr != nil {
				if err == nil {
					t.Errorf("fail %s: expecting error but got nil", test.name)
				}

				if !errors.Is(err, test.expErr) {
					t.Errorf("fail %s: expecting %q error but got : %q", test.name, test.expErr, err)
				}

				return
			}

			if got != test.amountCents {
				t.Errorf("Fail. Expecting %d but got %d", test.amountCents, got)
			}
		})
	}
}

func TestToAmount(t *testing.T) {
	tests := []struct {
		name        string
		amountCents int64
		amount      string
	}{
		{
			name:        "postitve amount",
			amountCents: 120,
			amount:      "1.20",
		},

		{
			name:        "negative amount",
			amountCents: -100,
			amount:      "-1.00",
		},

		{
			name:        "positive lesser than 1",
			amountCents: 1,
			amount:      "0.01",
		},

		{
			name:        "zero",
			amountCents: 0,
			amount:      "0.00",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ToAmount(test.amountCents)

			if got != test.amount {
				t.Errorf("Fail. Expected %s but got %s", test.amount, got)
			}
		})
	}
}
