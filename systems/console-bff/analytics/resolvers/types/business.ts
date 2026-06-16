/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Field, Float, InputType, Int, ObjectType } from "type-graphql";

import {
  ActivityItemDto,
  KpiDto,
  MetaDto,
  NamedValueDto,
  TimeSeriesDto,
} from "./shared";

@ObjectType()
export class SiteSummaryDto {
  @Field()
  siteId: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  status?: string;

  @Field(() => Float)
  revenue: number;

  @Field(() => Int)
  customers: number;
}

@ObjectType()
export class BusinessHomeDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [SiteSummaryDto])
  sites: SiteSummaryDto[];

  @Field(() => [NamedValueDto])
  topPackages: NamedValueDto[];

  @Field(() => [ActivityItemDto])
  recentActivity: ActivityItemDto[];
}

@ObjectType()
export class SalesOverviewDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => TimeSeriesDto, { nullable: true })
  revenueTrend?: TimeSeriesDto;

  @Field(() => [NamedValueDto])
  revenueBySite: NamedValueDto[];

  @Field(() => [NamedValueDto])
  revenueByPackage: NamedValueDto[];
}

@ObjectType()
export class PackageRowDto {
  @Field()
  packageId: string;

  @Field({ nullable: true })
  name?: string;

  @Field(() => Float)
  price: number;

  @Field({ nullable: true })
  validity?: string;

  @Field({ nullable: true })
  dataQuota?: string;

  @Field({ nullable: true })
  status?: string;

  @Field(() => Int)
  soldCount: number;

  @Field(() => Float)
  revenue: number;

  // This package's share of total revenue, 0–100. Backend gap: emitted empty
  // (null) for now; the console derives it client-side until the backend sends
  // it. See docs/analytics-backend-gaps.md.
  @Field(() => Float, { nullable: true })
  revenueSharePct?: number;

  @Field(() => Float)
  dataUsed: number;

  @Field(() => Int)
  activeSubscribers: number;
}

@ObjectType()
export class PackagePerformanceDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [PackageRowDto])
  packages: PackageRowDto[];

  @Field(() => [NamedValueDto])
  revenueMix: NamedValueDto[];

  @Field(() => MetaDto, { nullable: true })
  meta?: MetaDto;
}

@ObjectType()
export class InvoiceRowDto {
  @Field({ nullable: true })
  invoiceId?: string;

  @Field(() => Float)
  amount: number;

  @Field({ nullable: true })
  status?: string;

  @Field({ nullable: true })
  generatedAt?: string;
}

@ObjectType()
export class BillingSummaryDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => [InvoiceRowDto])
  invoices: InvoiceRowDto[];

  @Field({ nullable: true })
  lastInvoiceDate?: string;

  @Field(() => MetaDto, { nullable: true })
  meta?: MetaDto;
}

@ObjectType()
export class BusinessSiteRowDto {
  @Field()
  siteId: string;

  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  status?: string;

  @Field(() => Float)
  revenue: number;

  @Field(() => Float)
  revenueToday: number;

  @Field(() => Int)
  customers: number;

  @Field(() => Float)
  dataUsed: number;

  @Field(() => Float)
  uptime: number;

  @Field({ nullable: true })
  topPackage?: string;

  @Field({ nullable: true })
  issue?: string;

  @Field(() => Float)
  latitude: number;

  @Field(() => Float)
  longitude: number;
}

@ObjectType()
export class BusinessSitesDto {
  @Field(() => [BusinessSiteRowDto])
  sites: BusinessSiteRowDto[];

  @Field(() => MetaDto, { nullable: true })
  meta?: MetaDto;
}

@ObjectType()
export class BusinessSiteDto {
  @Field(() => BusinessSiteRowDto, { nullable: true })
  site?: BusinessSiteRowDto;

  @Field(() => [KpiDto])
  kpis: KpiDto[];

  @Field(() => TimeSeriesDto, { nullable: true })
  revenueTrend?: TimeSeriesDto;
}

@ObjectType()
export class InventoryReadinessDto {
  @Field(() => [KpiDto])
  kpis: KpiDto[];
}

@InputType()
export class AnalyticsSiteInput {
  @Field()
  siteId: string;

  @Field({ nullable: true })
  period?: string;

  @Field({ nullable: true })
  from?: string;

  @Field({ nullable: true })
  to?: string;

  @Field({ nullable: true })
  timezone?: string;
}
