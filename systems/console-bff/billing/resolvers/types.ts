/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class Owner {
  @Field()
  code: string;

  @Field()
  name: string;

  @Field()
  type: string;
}
@ObjectType()
export class Customer {
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
export class Fee {
  @Field()
  amountCents: number;

  @Field()
  amountCurrency: string;

  @Field()
  eventsCount: number;

  @Field(() => Owner)
  owner: Owner;

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
export class Fees {
  @Field(() => [Fee])
  fees: Fee[];
}

@ObjectType()
export class Subscription {
  @Field()
  canceldAt: string;

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

  @Field()
  terminatedAt: string;
}

@ObjectType()
export class Subscriptions {
  @Field(() => [Subscription])
  subscriptions: Subscription[];
}

@ObjectType()
export class GetReportsResDto {
  @Field(() => [GetReportResDto])
  reports: GetReportResDto[];
}

@InputType()
export class GetReportsInputDto {
  @Field()
  networkId?: string;

  @Field()
  ownerId?: string;

  @Field()
  ownerType?: string;

  @Field()
  report_type?: string;

  @Field()
  isPaid?: boolean;

  @Field()
  sort?: boolean;

  @Field()
  count?: number;
}

@ObjectType()
export class RawReport {
  @Field()
  amountCents: number;

  @Field()
  amountCurrency: string;

  @Field(() => Customer)
  customer: Customer;

  @Field(() => Fees)
  fees: Fees;

  @Field(() => Subscriptions)
  subscriptions: Subscriptions;

  @Field()
  fileURL: string;

  @Field()
  issuingDate: string;

  @Field()
  paymentStatus: string;

  @Field()
  status: string;

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
export class GetReportResDto {
  @Field()
  type: string;

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

  @Field(() => RawReport)
  rawReport: RawReport;
}
