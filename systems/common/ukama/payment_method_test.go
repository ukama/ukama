/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
)

func TestPaymentMethod(t *testing.T) {
	t.Run("PaymentMethodValidString", func(tt *testing.T) {
		stripeMethod := ukama.ParsePaymentMethod("StrIpE")

		assert.NotNil(t, stripeMethod)
		assert.Equal(t, stripeMethod.String(), ukama.PaymentMethodStripe.String())
		assert.Equal(t, uint8(stripeMethod), uint8(1))
	})

	t.Run("PaymentMethodValidNumber", func(tt *testing.T) {
		mopayMethod := ukama.ParsePaymentMethod("2")

		assert.NotNil(t, mopayMethod)
		assert.Equal(t, uint8(mopayMethod), uint8(2))
		assert.Equal(t, mopayMethod.String(), ukama.PaymentMethodMoPay.String())
	})

	t.Run("PaymentMethodNonValidString", func(tt *testing.T) {
		unsupportedMethod := ukama.ParsePaymentMethod("failure")

		assert.NotNil(t, unsupportedMethod)
		assert.Equal(t, unsupportedMethod.String(), ukama.PaymentMethodUnknown.String())
		assert.Equal(t, uint8(unsupportedMethod), uint8(0))
	})

	t.Run("PaymentMethodNonValidNumber", func(tt *testing.T) {
		unknownMethod := ukama.PaymentMethod(uint8(10))

		assert.NotNil(t, unknownMethod)
		assert.Equal(t, unknownMethod.String(), ukama.PaymentMethodUnknown.String())
		assert.Equal(t, uint8(unknownMethod), uint8(10))
	})
}
