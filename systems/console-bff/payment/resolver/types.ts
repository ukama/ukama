/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

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
export class UpdatePaymentInputDto {
  @Field()
  id: string;

  @Field({ nullable: true })
  country?: string;

  @Field({ nullable: true })
  currency?: string;

  @Field({ nullable: true })
  payerEmail?: string;

  @Field({ nullable: true })
  payerName?: string;

  @Field({ nullable: true })
  paymentMethod?: string;
}

@InputType()
export class ProcessPaymentInputDto {
  @Field()
  id: string;

  @Field({ nullable: true })
  correspondent?: string;

  @Field()
  token: string;
}

@ObjectType()
export class TokenResDto {
  @Field()
  token: string;
}

@InputType()
export class GetPaymentsInputDto {
  @Field({ nullable: true })
  type?: string;

  @Field({ nullable: true })
  paymentMethod?: string;

  @Field({ nullable: true })
  status?: string;
}
