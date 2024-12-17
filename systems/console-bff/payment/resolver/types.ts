/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

import { PAYMENT_ITEM_TYPE } from "../../common/enums";

@ObjectType()
export class PaymentDto {
  @Field()
  id: string;

  @Field()
  itemId: string;

  @Field()
  itemType: string;

  @Field()
  amount: string;

  @Field()
  currency: string;

  @Field()
  paymentMethod: string;

  @Field()
  depositedAmount: string;

  @Field()
  paidAt: string;

  @Field()
  payerName: string;

  @Field()
  payerEmail: string;

  @Field()
  payerPhone: string;

  @Field()
  correspondent: string;

  @Field()
  country: string;

  @Field()
  description: string;

  @Field()
  status: string;

  @Field()
  extra: string;

  @Field()
  failureReason: string;

  @Field()
  createdAt: string;
}

@ObjectType()
export class ProcessPaymentDto {
  @Field(() => PaymentDto)
  payment: PaymentDto;
}

@ObjectType()
export class PaymentsDto {
  @Field(() => [PaymentDto])
  payments: PaymentDto[];
}

@ObjectType()
export class PaymentAPIDto {
  @Field()
  id: string;

  @Field()
  item_id: string;

  @Field()
  item_type: string;

  @Field()
  amount: string;

  @Field()
  currency: string;

  @Field()
  payment_method: string;

  @Field()
  deposited_amount: string;

  @Field()
  paid_at: string;

  @Field()
  payer_name: string;

  @Field()
  payer_email: string;

  @Field()
  payer_phone: string;

  @Field()
  correspondent: string;

  @Field()
  country: string;

  @Field()
  description: string;

  @Field()
  status: string;

  @Field()
  extra: string;

  @Field()
  failure_reason: string;

  @Field()
  created_at: string;
}

@ObjectType()
export class ProcessPaymentAPIResDto {
  @Field(() => PaymentAPIDto)
  payment: PaymentAPIDto;
}

@ObjectType()
export class PaymentAPIResDto {
  @Field(() => PaymentAPIDto)
  payment: PaymentAPIDto;
}

@ObjectType()
export class PaymentsAPIResDto {
  @Field(() => [PaymentAPIDto])
  payments: PaymentAPIDto[];
}

@InputType()
export class AddPaymentInputDto {
  @Field()
  amount: string;

  @Field()
  country: string;

  @Field()
  currency: string;

  @Field()
  description: string;

  @Field()
  itemId: string;

  @Field(() => PAYMENT_ITEM_TYPE)
  itemType: PAYMENT_ITEM_TYPE;

  @Field()
  payerEmail: string;

  @Field()
  payerPhone: string;

  @Field()
  paymentMethod: string;

  @Field()
  correspondent?: string;

  @Field()
  simId: string;

  @Field()
  subscriberId: string;
}

@InputType()
export class UpdatePaymentInputDto {
  @Field()
  id: string;

  @Field()
  country: string;

  @Field()
  currency: string;

  @Field()
  description: string;

  @Field()
  amount: string;

  @Field()
  paidAt: string;

  @Field()
  payerEmail: string;

  @Field()
  payerName: string;

  @Field()
  payerPhone: string;

  @Field()
  paymentMethod: string;

  @Field()
  status: string;
}

@InputType()
export class ProcessPaymentInputDto {
  @Field()
  id: string;

  @Field()
  correspondent: string;

  @Field()
  token: string;
}

@InputType()
export class CorrespondentsInputDto {
  @Field()
  phoneNumber: string;

  @Field()
  paymentMethod: string;
}

@ObjectType()
export class CorrespondentsAPIDto {
  @Field(() => [String])
  correspondents: string[];

  @Field()
  country: string;
}

@ObjectType()
export class CorrespondentsDto {
  @Field()
  label: string;

  @Field()
  logo: string;

  @Field()
  correspondent_code: string;
}

@ObjectType()
export class CorrespondentsResDto {
  @Field(() => [CorrespondentsDto])
  correspondents: CorrespondentsDto[];

  @Field()
  country: string;
}

@ObjectType()
export class TokenResDto {
  @Field()
  token: string;
}

@InputType()
export class GetPaymentsInputDto {
  @Field()
  type: string;

  @Field()
  paymentMethod: string;

  @Field()
  status: string;
}
