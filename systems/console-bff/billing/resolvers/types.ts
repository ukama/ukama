/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class ItemResDto {
  @Field()
  type: string;

  @Field()
  code: string;

  @Field()
  name: string;
}

@InputType()
export class GetReportsInputDto {
  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  ownerId?: string;

  @Field()
  ownerType?: string;

  @Field()
  report_type?: string;

  @Field({ nullable: true })
  isPaid?: boolean;

  @Field({ nullable: true })
  sort?: boolean;

  @Field({ nullable: true })
  count?: number;
}

@ObjectType()
export class CustomerResDto {
  @Field()
  external_id: string;

  @Field()
  name: string;

  @Field({ nullable: true })
  email: string;

  @Field({ nullable: true })
  address_line1: string;

  @Field({ nullable: true })
  legal_name: string;

  @Field({ nullable: true })
  legal_number: string;

  @Field({ nullable: true })
  phone: string;

  @Field()
  currency: string;

  @Field({ nullable: true })
  timezone: string;

  @Field()
  vat_rate: number;

  @Field()
  created_at: string;
}

@ObjectType()
export class FeeResDto {
  @Field()
  taxes_amount_cents: string;

  @Field()
  taxes_precise_amount: string;

  @Field()
  total_amount_cents: string;

  @Field()
  total_amount_currency: string;

  @Field()
  events_count: string;

  @Field()
  units: number;

  @Field(() => ItemResDto)
  item: ItemResDto;
}

@ObjectType()
export class SubscriptionResDto {
  @Field()
  external_customer_id: string;

  @Field()
  external_id: string;

  @Field()
  plan_code: string;

  @Field({ nullable: true })
  name: string;

  @Field()
  status: string;

  @Field()
  created_at: string;

  @Field()
  started_at: string;

  @Field({ nullable: true })
  canceled_at: string;

  @Field({ nullable: true })
  terminated_at: string;
}

@ObjectType()
export class RawReportResDto {
  @Field()
  issuing_date: string;

  @Field()
  payment_due_date: string;

  @Field()
  payment_overdue: boolean;

  @Field()
  invoice_type: string;

  @Field()
  status: string;

  @Field()
  payment_status: string;

  @Field()
  fees_amount_cents: string;

  @Field()
  taxes_amount_cents: string;

  @Field()
  sub_total_excluding_taxes_amount_cents: string;

  @Field()
  sub_total_including_taxes_amount_cents: string;

  @Field()
  vat_amount_cents: string;

  @Field({ nullable: true })
  vat_amount_currency: string;

  @Field()
  total_amount_cents: string;

  @Field()
  currency: string;

  @Field()
  file_url: string;

  @Field(() => CustomerResDto)
  customer: CustomerResDto;

  @Field(() => [SubscriptionResDto])
  subscriptions: SubscriptionResDto[];

  @Field(() => [FeeResDto])
  fees: FeeResDto[];
}

@ObjectType()
export class ReportResDto {
  @Field()
  id: string;

  @Field()
  owner_id: string;

  @Field()
  owner_type: string;

  @Field()
  network_Id: string;

  @Field()
  period: string;

  @Field()
  Type: string;

  @Field(() => RawReportResDto)
  raw_report: RawReportResDto;

  @Field()
  is_paid: boolean;

  @Field()
  created_at: string;
}

@ObjectType()
export class CustomerDto {
  @Field()
  externalId: string;

  @Field()
  name: string;

  @Field({ nullable: true })
  email: string;

  @Field({ nullable: true })
  addressLine1: string;

  @Field({ nullable: true })
  legalName: string;

  @Field({ nullable: true })
  legalNumber: string;

  @Field({ nullable: true })
  phone: string;

  @Field()
  currency: string;

  @Field({ nullable: true })
  timezone: string;

  @Field()
  vatRate: number;

  @Field()
  createdAt: string;
}

@ObjectType()
export class FeeDto {
  @Field()
  taxesAmountCents: string;

  @Field()
  taxesPreciseAmount: string;

  @Field()
  totalAmountCents: string;

  @Field()
  totalAmountCurrency: string;

  @Field()
  eventsCount: string;

  @Field()
  units: number;

  @Field(() => ItemResDto)
  item: ItemResDto;
}

@ObjectType()
export class SubscriptionDto {
  @Field()
  externalCustomerId: string;

  @Field()
  externalId: string;

  @Field()
  planCode: string;

  @Field({ nullable: true })
  name: string;

  @Field()
  status: string;

  @Field()
  createdAt: string;

  @Field()
  startedAt: string;

  @Field({ nullable: true })
  canceledAt: string;

  @Field({ nullable: true })
  terminatedAt: string;
}

@ObjectType()
export class RawReportDto {
  @Field()
  issuingDate: string;

  @Field()
  paymentDueDate: string;

  @Field()
  paymentOverdue: boolean;

  @Field()
  invoiceType: string;

  @Field()
  status: string;

  @Field()
  paymentStatus: string;

  @Field()
  feesAmountCents: string;

  @Field()
  taxesAmountCents: string;

  @Field()
  subTotalExcludingTaxesAmountCents: string;

  @Field()
  subTotalIncludingTaxesAmountCents: string;

  @Field()
  vatAmountCents: string;

  @Field({ nullable: true })
  vatAmountCurrency: string;

  @Field()
  totalAmountCents: string;

  @Field()
  currency: string;

  @Field()
  fileUrl: string;

  @Field(() => CustomerDto)
  customer: CustomerDto;

  @Field(() => [SubscriptionDto])
  subscriptions: SubscriptionDto[];

  @Field(() => [FeeDto])
  fees: FeeDto[];
}

@ObjectType()
export class ReportDto {
  @Field()
  id: string;

  @Field()
  ownerId: string;

  @Field()
  ownerType: string;

  @Field()
  networkId: string;

  @Field()
  period: string;

  @Field()
  type: string;

  @Field(() => RawReportDto)
  rawReport: RawReportDto;

  @Field()
  isPaid: boolean;

  @Field()
  createdAt: string;
}

@ObjectType()
export class GetReportsDto {
  @Field(() => [ReportDto])
  reports: ReportDto[];
}

@ObjectType()
export class GetReportsResDto {
  @Field(() => [ReportResDto])
  reports: ReportResDto[];
}

@ObjectType()
export class GetReportDto {
  @Field(() => ReportDto)
  report: ReportDto;
}

@ObjectType()
export class GetReportResDto {
  @Field(() => ReportResDto)
  report: ReportResDto;
}
