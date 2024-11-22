/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

type PaymentMethod uint8

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodStripe
	PaymentMethodMoPay
)

func (s *PaymentMethod) Scan(value interface{}) error {
	*s = PaymentMethod(uint8(value.(int64)))
	return nil
}

func (s PaymentMethod) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s PaymentMethod) String() string {
	t := map[PaymentMethod]string{0: "unknown", 1: "stripe", 2: "mopay"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParsePaymentMethod(value string) PaymentMethod {
	i, err := strconv.Atoi(value)
	if err == nil {
		return PaymentMethod(i)
	}

	t := map[string]PaymentMethod{"unknown": 0, "stripe": 1, "mopay": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return PaymentMethod(0)
	}

	return PaymentMethod(v)
}
