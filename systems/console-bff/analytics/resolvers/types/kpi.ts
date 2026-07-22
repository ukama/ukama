/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Field, Float, InputType, ObjectType } from "type-graphql";

/**
 * GraphQL types for the analytics gateway's generic KPI read API
 * (`GET /v1/analytics/kpis/values`). The gateway renders protobuf responses
 * with protojson (EmitUnpopulated), so the wire JSON already uses
 * lowerCamelCase keys — these fields map 1:1 to the `KpiValue`/`Trend`
 * messages. The only reshape is `scope` (a protobuf map) → `[ScopeEntryDto]`,
 * done in the datasource so it is expressible in GraphQL.
 */

/** Day/week/month-over-previous comparison for a KPI value. */
@ObjectType()
export class TrendDto {
  @Field({ nullable: true })
  direction?: string; // up | down | flat | new | na

  @Field(() => Float, { nullable: true })
  changePct?: number;

  @Field(() => Float, { nullable: true })
  changeAbs?: number;

  @Field(() => Float, { nullable: true })
  prevValue?: number;

  @Field({ nullable: true })
  hasPrevious?: boolean;
}

/** One entry of a KPI value's scope map, e.g. { key: "network_id", value }. */
@ObjectType()
export class ScopeEntryDto {
  @Field()
  key: string;

  @Field()
  value: string;
}

/**
 * A single KPI reading: the metadata envelope of `<kpi>: <value>` plus the
 * span/op it was rolled up at, the window it covers, unit/symbol/type and the
 * trend vs. the previous window.
 */
@ObjectType()
export class KpiValueDto {
  @Field()
  kpi: string;

  @Field(() => Float)
  value: number;

  @Field({ nullable: true })
  span?: string; // daily | weekly | monthly

  @Field({ nullable: true })
  op?: string;

  @Field({ nullable: true })
  from?: string; // RFC3339

  @Field({ nullable: true })
  to?: string; // RFC3339

  @Field({ nullable: true })
  type?: string;

  @Field({ nullable: true })
  unit?: string;

  @Field({ nullable: true })
  symbol?: string;

  @Field({ nullable: true })
  isPartial?: boolean;

  @Field(() => [ScopeEntryDto])
  scope: ScopeEntryDto[];

  @Field(() => TrendDto, { nullable: true })
  trend?: TrendDto;

  @Field({ nullable: true })
  computedAt?: string; // RFC3339; empty when no data yet
}

/** Response envelope for `getKpiValues` (mirrors GetKpisResponse). */
@ObjectType()
export class GetKpiValuesDto {
  @Field(() => [KpiValueDto])
  values: KpiValueDto[];
}

/** Response envelope for `getKpiTimeSeries` — one KpiValue per span bucket in
 *  the range (mirrors GetKpiTimeSeriesResponse). */
@ObjectType()
export class GetKpiTimeSeriesDto {
  @Field(() => [KpiValueDto])
  values: KpiValueDto[];
}

/**
 * Input for `getKpiValues`. `keys` is required (the KPI keys to read, e.g.
 * ["network_uptime", "site_uptime"]); `span`/`op` default per the KPI spec on
 * the gateway; `networkId` is an optional scope filter.
 */
@InputType()
export class KpiValuesInput {
  @Field(() => [String])
  keys: string[];

  @Field({ nullable: true })
  span?: string; // daily | weekly | monthly (default daily)

  @Field({ nullable: true })
  op?: string; // AVG | MIN | MAX | ... (defaults per KPI spec)

  @Field({ nullable: true })
  networkId?: string;
}

/**
 * Input for `getKpiTimeSeries`: one value per span bucket over a range. Same as
 * KpiValuesInput plus the [from, to) window (RFC3339). `span` is the bucket
 * size (e.g. weekly for a 9-week revenue trend); defaults per KPI spec.
 */
@InputType()
export class KpiTimeSeriesInput {
  @Field(() => [String])
  keys: string[];

  @Field({ nullable: true })
  span?: string; // daily | weekly | monthly

  @Field({ nullable: true })
  op?: string;

  @Field({ nullable: true })
  from?: string; // RFC3339 inclusive

  @Field({ nullable: true })
  to?: string; // RFC3339 exclusive

  @Field({ nullable: true })
  networkId?: string;
}
