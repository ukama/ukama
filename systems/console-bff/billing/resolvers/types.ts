/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

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

@InputType()
export class GetReportsInputDto {
  @Field({ nullable: true })
  count?: number;

  @Field({ nullable: true })
  is_paid?: boolean;

  @Field()
  network_id: string;

  @Field({ nullable: true })
  owner_id?: string;

  @Field({ nullable: true })
  owner_type?: string;

  @Field({ nullable: true })
  report_type?: string;
}
@InputType()
export class GetReportInputDto {
  @Field()
  id: string;

  @Field()
  asPdf: boolean;
}

@ObjectType()
class CustomerDto {
  @Field()
  AddressLine1: string;

  @Field()
  createdAt: string;

  @Field()
  email: string;

  @Field()
  externalId: string;

  @Field()
  legalName: string;

  @Field()
  legalNumber: string;

  @Field()
  name: string;

  @Field()
  phone: string;

  @Field()
  vatRate: number;
}

@ObjectType()
class FeeOwnerDto {
  @Field()
  code: string;

  @Field()
  name: string;

  @Field()
  type: string;
}

@ObjectType()
class FeeDto {
  @Field()
  amountCents: number;

  @Field()
  amountCurrency: string;

  @Field()
  eventsCount: number;

  @Field(() => FeeOwnerDto)
  owner: FeeOwnerDto;

  @Field()
  totalAmountCents: number;

  @Field()
  totalAmountCurrency: string;

  @Field()
  units: number;

  @Field()
  vatAmountCents: number;

  @Field()
  vatAmountCurrency: string;
}

@ObjectType()
class SubscriptionDto {
  @Field({ nullable: true })
  canceldAt?: string;

  @Field()
  createdAt: string;

  @Field()
  externalCustomerId: string;

  @Field()
  externalId: string;

  @Field()
  planCode: string;

  @Field()
  startedAt: string;

  @Field()
  status: string;

  @Field({ nullable: true })
  terminatedAt?: string;
}

@ObjectType()
class RawReportDto {
  @Field()
  amountCents: number;

  @Field()
  amountCurrency: string;

  @Field(() => CustomerDto)
  customer: CustomerDto;

  @Field(() => [FeeDto])
  fees: FeeDto[];

  @Field()
  fileURL: string;

  @Field()
  issuingDate: string;

  @Field()
  paymentStatus: string;

  @Field()
  status: string;

  @Field(() => [SubscriptionDto])
  subscriptions: SubscriptionDto[];

  @Field()
  totalAmountCents: number;

  @Field()
  totalAmountCurrency: string;

  @Field()
  vatAmountCents: number;

  @Field()
  vatAmountCurrency: string;
}

@ObjectType()
class ReportDto {
  @Field()
  Type: string;

  @Field()
  createdAt: string;

  @Field()
  id: string;

  @Field()
  isPaid: boolean;

  @Field()
  networkId: string;

  @Field()
  ownerId: string;

  @Field()
  ownerType: string;

  @Field()
  period: string;

  @Field(() => RawReportDto)
  rawReport: RawReportDto;
}

@ObjectType()
export class GetReportsResDto {
  @Field(() => [ReportDto])
  reports: ReportDto[];
}

@ObjectType()
export class GetReportResDto {
  @Field()
  report: ReportDto;
}

@InputType()
export class AddReportInputDto {
  @Field(() => InvoiceInputDto)
  invoice: InvoiceInputDto;

  @Field()
  object_type: string;

  @Field()
  webhook_type: string;
}

@InputType()
export class InvoiceInputDto {
  @Field({ nullable: true })
  amount_cents?: number;

  @Field({ nullable: true })
  amount_currency?: string;

  @Field({ nullable: true })
  credit_amount_cents?: number;

  @Field({ nullable: true })
  credit_amount_currency?: string;

  @Field(() => [CreditInputDto], { nullable: true })
  credits?: CreditInputDto[];

  @Field(() => CustomerInputDto)
  customer: CustomerInputDto;

  @Field(() => [FeeInputDto])
  fees: FeeInputDto[];

  @Field({ nullable: true })
  file_url?: string;

  @Field({ nullable: true })
  issuing_date?: string;

  @Field({ nullable: true })
  legacy?: boolean;

  @Field(() => [MetadataInputDto], { nullable: true })
  metadata?: MetadataInputDto[];

  @Field({ nullable: true })
  number?: string;

  @Field({ nullable: true })
  payment_status?: string;

  @Field({ nullable: true })
  sequential_id?: number;

  @Field({ nullable: true })
  status?: string;

  @Field(() => [SubscriptionInputDto], { nullable: true })
  subscriptions?: SubscriptionInputDto[];

  @Field({ nullable: true })
  total_amount_cents?: number;

  @Field({ nullable: true })
  total_amount_currency?: string;

  @Field({ nullable: true })
  vat_amount_cents?: number;

  @Field({ nullable: true })
  vat_amount_currency?: string;
}

@InputType()
export class CreditInputDto {
  @Field({ nullable: true })
  amount_cents?: number;

  @Field({ nullable: true })
  amount_currency?: string;

  @Field(() => CreditItemInputDto)
  item: CreditItemInputDto;

  @Field({ nullable: true })
  lago_id?: string;
}

@InputType()
export class CreditItemInputDto {
  @Field({ nullable: true })
  code?: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  type?: string;
}

@InputType()
export class CustomerInputDto {
  @Field({ nullable: true })
  address_line1?: string;

  @Field({ nullable: true })
  address_line2?: string;

  @Field({ nullable: true })
  city?: string;

  @Field({ nullable: true })
  country?: string;

  @Field({ nullable: true })
  created_at?: string;

  @Field({ nullable: true })
  email?: string;

  @Field({ nullable: true })
  external_id?: string;

  @Field({ nullable: true })
  legal_name?: string;

  @Field({ nullable: true })
  legal_number?: string;

  @Field({ nullable: true })
  logo_url?: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  phone?: string;

  @Field({ nullable: true })
  state?: string;

  @Field({ nullable: true })
  url?: string;

  @Field({ nullable: true })
  vat_rate?: number;

  @Field({ nullable: true })
  zipcode?: string;
}

@InputType()
export class FeeInputDto {
  @Field({ nullable: true })
  amount_cents?: number;

  @Field({ nullable: true })
  amount_currenty?: string;

  @Field({ nullable: true })
  events_count?: number;

  @Field(() => FeeItemInputDto)
  item: FeeItemInputDto;

  @Field({ nullable: true })
  total_amount_cents?: number;

  @Field({ nullable: true })
  total_amount_currency?: string;

  @Field({ nullable: true })
  units?: string;

  @Field({ nullable: true })
  vat_amount_cents?: number;

  @Field({ nullable: true })
  vat_amount_currency?: string;
}

@InputType()
export class FeeItemInputDto {
  @Field({ nullable: true })
  code?: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  type?: string;
}

@InputType()
export class MetadataInputDto {
  @Field({ nullable: true })
  created_at?: string;

  @Field({ nullable: true })
  key?: string;

  @Field({ nullable: true })
  value?: string;
}

@InputType()
export class SubscriptionInputDto {
  @Field({ nullable: true })
  canceled_at?: string;

  @Field({ nullable: true })
  created_at?: string;

  @Field({ nullable: true })
  external_customer_id?: string;

  @Field({ nullable: true })
  external_id?: string;

  @Field({ nullable: true })
  plan_code?: string;

  @Field({ nullable: true })
  started_at?: string;

  @Field({ nullable: true })
  status?: string;

  @Field({ nullable: true })
  terminated_at?: string;
}
