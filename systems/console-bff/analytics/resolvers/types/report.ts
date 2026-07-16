/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Field, Float, InputType, Int, ObjectType } from "type-graphql";

import { ScopeEntryDto, TrendDto } from "./kpi";

/**
 * GraphQL types for the analytics gateway's performance report API
 * (`GET /v1/analytics/reports/{report}`), which composes a resource
 * performance table from the latest available KPI values plus entity
 * attributes. Wire shape is protojson (lowerCamelCase); the only reshape is
 * each row's `attributes` map → `[ScopeEntryDto]`, done in the datasource.
 */

/** One computed cell (column) of a report row. */
@ObjectType()
export class ReportCellDto {
  @Field()
  column: string;

  @Field(() => Float)
  value: number;

  @Field({ nullable: true })
  unit?: string;

  @Field({ nullable: true })
  symbol?: string;

  @Field({ nullable: true })
  format?: string; // display hint from the report spec

  @Field(() => TrendDto, { nullable: true })
  trend?: TrendDto;

  @Field({ nullable: true })
  isPartial?: boolean;

  @Field({ nullable: true })
  computedAt?: string; // RFC3339; empty when no data yet
}

/** One entity's row: its id, attributes, computed cells and status. */
@ObjectType()
export class ReportRowDto {
  @Field()
  entityId: string;

  @Field(() => [ScopeEntryDto])
  attributes: ScopeEntryDto[];

  @Field(() => [ReportCellDto])
  cells: ReportCellDto[];

  @Field({ nullable: true })
  status?: string;
}

/** Response envelope for `getPerformanceReport`. */
@ObjectType()
export class GetPerformanceReportDto {
  @Field()
  report: string;

  @Field({ nullable: true })
  title?: string;

  @Field({ nullable: true })
  span?: string;

  @Field(() => [ReportRowDto])
  rows: ReportRowDto[];
}

/**
 * Input for `getPerformanceReport`. `report` is the report key from the
 * registry (`GET /reports`); `span`/`networkId` scope the table; `top` caps
 * the number of rows.
 */
@InputType()
export class PerformanceReportInput {
  @Field()
  report: string;

  @Field({ nullable: true })
  span?: string; // daily | weekly | monthly (default daily)

  @Field({ nullable: true })
  networkId?: string;

  @Field(() => Int, { nullable: true })
  top?: number;
}
