package util

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
			name:        "mal formed positive amount",
			amount:      "1.20",
			amountCents: 0,
			expErr:      ErrMatch,
		},

		{
			name:        "negative amount",
			amount:      "-1.5",
			amountCents: 0,
			expErr:      ErrMatch,
		},

		{
			name:        "positive lesser than 100",
			amount:      "0.1",
			amountCents: 10,
			expErr:      nil,
		},

		{
			name:        "zero",
			amount:      "0.00",
			amountCents: 0,
			expErr:      ErrMatch,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ToAmountCents(test.amount)

			if test.expErr != nil {
				if err == nil {
					t.Errorf("fail, expecting error but got nil")
				}

				if !errors.Is(err, test.expErr) {
					t.Errorf("fail, expecting %q error but got : %q", test.expErr, err)
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
			name:        "positive lesser than 100",
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
