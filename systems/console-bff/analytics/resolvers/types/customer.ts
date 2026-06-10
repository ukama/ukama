/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Field, Float, InputType, Int, ObjectType } from "type-graphql";

import { ActivityItemDto, KpiDto, MetaDto } from "./shared";

@ObjectType()
export class CustomerRowDto {
  @Field()
  customerId: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  email?: string;

  @Field({ nullable: true })
  status?: string;

  @Field({ nullable: true })
  packageName?: string;

  @Field({ nullable: true })
  packageStatus?: string;

  @Field(() => Float)
  dataUsage: number;

  @Field({ nullable: true })
  lastSeen?: string;

  @Field({ nullable: true })
  siteId?: string;

  @Field({ nullable: true })
  siteName?: string;

  @Field({ nullable: true })
  simIccid?: string;

  @Field({ nullable: true })
  simStatus?: string;
}

@ObjectType()
export class CustomerOverviewDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];
}

@ObjectType()
export class CustomerListDto {
  @Field(() => [CustomerRowDto])
  customers: CustomerRowDto[];

  @Field(() => MetaDto, { nullable: true })
  meta?: MetaDto;
}

@ObjectType()
export class PackageIntervalDto {
  @Field({ nullable: true })
  packageId?: string;

  @Field({ nullable: true })
  packageName?: string;

  @Field({ nullable: true })
  startAt?: string;

  @Field({ nullable: true })
  endAt?: string;

  @Field({ nullable: true })
  state?: string;
}

@ObjectType()
export class CustomerDetailDto {
  @Field(() => CustomerRowDto, { nullable: true })
  customer?: CustomerRowDto;

  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [PackageIntervalDto])
  packageHistory: PackageIntervalDto[];
}

@ObjectType()
export class SupportSignalDto {
  @Field()
  key: string;

  @Field({ nullable: true })
  state?: string;

  @Field({ nullable: true })
  detail?: string;
}

@ObjectType()
export class CustomerSupportDto {
  @Field(() => CustomerRowDto, { nullable: true })
  customer?: CustomerRowDto;

  @Field({ nullable: true })
  likelyIssue?: string;

  @Field({ nullable: true })
  recommendedAction?: string;

  @Field()
  escalationNeeded: boolean;

  @Field(() => [SupportSignalDto])
  signals: SupportSignalDto[];

  @Field(() => [ActivityItemDto])
  recentActivity: ActivityItemDto[];
}

@ObjectType()
export class SimRowDto {
  @Field()
  simId: string;

  @Field({ nullable: true })
  iccid?: string;

  @Field({ nullable: true })
  status?: string;

  @Field({ nullable: true })
  customerId?: string;

  @Field({ nullable: true })
  batchId?: string;

  @Field({ nullable: true })
  allocatedAt?: string;
}

@ObjectType()
export class CustomerSimsDto {
  @Field(() => [SimRowDto])
  sims: SimRowDto[];

  @Field(() => MetaDto, { nullable: true })
  meta?: MetaDto;
}

@ObjectType()
export class SimBatchDto {
  @Field()
  batchId: string;

  @Field(() => Int)
  quantity: number;

  @Field(() => Int)
  assigned: number;

  @Field(() => Float)
  assignedPercent: number;

  @Field({ nullable: true })
  uploadedAt?: string;
}

@ObjectType()
export class SimPoolDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [SimBatchDto])
  batches: SimBatchDto[];
}

@InputType()
export class CustomerByIdInput {
  @Field()
  customerId: string;

  @Field({ nullable: true })
  period?: string;

  @Field({ nullable: true })
  from?: string;

  @Field({ nullable: true })
  to?: string;

  @Field({ nullable: true })
  timezone?: string;
}
