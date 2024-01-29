/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, Float, InputType, Int, ObjectType } from "type-graphql";

@ObjectType()
export class CurrentBillDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  dataUsed: number;

  @Field()
  rate: number;

  @Field()
  subtotal: number;
}

@ObjectType()
export class CurrentBillResponse {
  @Field()
  status: string;

  @Field(() => [CurrentBillDto])
  data: CurrentBillDto[];
}

@ObjectType()
export class BillResponse {
  @Field(() => [CurrentBillDto])
  bill: CurrentBillDto[];

  @Field()
  total: number;

  @Field()
  billMonth: string;

  @Field()
  dueDate: string;
}

@ObjectType()
export class BillHistoryDto {
  @Field()
  id: string;

  @Field()
  date: string;

  @Field()
  description: string;

  @Field()
  totalUsage: number;

  @Field()
  subtotal: number;
}

@ObjectType()
export class BillHistoryResponse {
  @Field()
  status: string;

  @Field(() => [BillHistoryDto])
  data: BillHistoryDto[];
}
@ObjectType()
export class StripeCustomer {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  email: string;
}

@InputType()
export class CreateCustomerDto {
  @Field()
  name: string;

  @Field()
  email: string;
}

@ObjectType()
export class Fee {
  @Field(() => Int)
  amountCents: number;

  @Field()
  amountCurrency: string;

  @Field(() => Int)
  vatAmountCents: number;

  @Field()
  vatAmountCurrency: string;

  @Field(() => Int)
  totalAmountCents: number;

  @Field()
  totalAmountCurrency: string;

  @Field(() => Int)
  eventsCount: number;

  @Field(() => Float)
  units: number;

  @Field(() => FeeItem)
  item: FeeItem;
}

@ObjectType()
export class FeeItem {
  @Field()
  type: string;

  @Field()
  code: string;

  @Field()
  name: string;
}
@ObjectType()
export class RawInvoiceDto {
  @Field()
  issuingDate: string;

  @Field()
  status: string;

  @Field()
  paymentStatus: string;

  @Field(() => Int)
  amountCents: number;

  @Field()
  amountCurrency: string;

  @Field(() => Int)
  vatAmountCents: number;

  @Field()
  vatAmountCurrency: string;

  @Field(() => Int)
  totalAmountCents: number;

  @Field()
  totalAmountCurrency: string;

  @Field()
  fileURL: string;

  // Assuming Customer, Subscription, and Fee are also GraphQL types
  @Field(() => Customer)
  customer: Customer;

  @Field(() => [Subscription])
  subscriptions: Subscription[];

  @Field(() => [Fee])
  fees: Fee[];
}
export class Subscription {
  @Field()
  externalCustomerId: string;

  @Field()
  externalId: string;

  @Field()
  planCode: string;

  @Field()
  status: string;

  @Field()
  createdAt: string;

  @Field()
  startedAt: string;

  @Field()
  canceledAt: string;

  @Field()
  terminatedAt: string;
}
@ObjectType()
export class Customer {
  @Field()
  externalId: string;

  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  addressLine1: string;

  @Field()
  legalName: string;

  @Field()
  legalNumber: string;

  @Field()
  phone: string;

  @Field(() => Float)
  vatRate: number;

  @Field()
  createdAt: string;
}
@ObjectType()
export class InvoiceDto {
  @Field()
  id: string;

  @Field()
  susbcriberId: string;

  @Field()
  networkId: string;

  @Field()
  period: Date;

  @Field()
  rawInvoice: string;

  @Field()
  isPaid: boolean;

  @Field()
  createdAt: Date;
}

@ObjectType()
export class StripePaymentMethods {
  @Field()
  id: string;

  @Field()
  brand: string;

  @Field({ nullable: true })
  cvc_check?: string;

  @Field({ nullable: true })
  country?: string;

  @Field()
  exp_month: number;

  @Field()
  exp_year: number;

  @Field()
  funding: string;

  @Field()
  last4: string;

  @Field()
  type: string;

  @Field()
  created: number;
}
