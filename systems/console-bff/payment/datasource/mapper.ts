/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  PaymentAPIDto,
  PaymentAPIResDto,
  PaymentDto,
  PaymentsAPIResDto,
  PaymentsDto,
  ProcessPaymentAPIResDto,
  ProcessPaymentDto,
} from "../resolver/types";

export const paymentDtoMapper = (req: PaymentAPIDto): PaymentDto => {
  return {
    id: req.id,
    itemId: req.item_id,
    itemType: req.item_type,
    amount: req.amount,
    currency: req.currency,
    paymentMethod: req.payment_method,
    depositedAmount: req.deposited_amount,
    paidAt: req.paid_at,
    payerName: req.payer_name,
    payerEmail: req.payer_email,
    payerPhone: req.payer_phone,
    correspondent: req.correspondent,
    country: req.country,
    description: req.description,
    status: req.status,
    failureReason: req.failure_reason,
    extra: req.extra,
    createdAt: req.created_at,
  };
};

export const dtoToPaymentDto = (res: PaymentAPIResDto): PaymentDto => {
  return paymentDtoMapper(res.payment);
};

export const dtoToPaymentsDto = (res: PaymentsAPIResDto): PaymentsDto => {
  const payments: PaymentDto[] = [];
  res.payments.forEach(payment => {
    payments.push(paymentDtoMapper(payment));
  });
  return {
    payments,
  };
};

export const dtoToProcessPaymentsDto = (
  res: ProcessPaymentAPIResDto
): ProcessPaymentDto => {
  return {
    payment: paymentDtoMapper(res.payment),
  };
};
