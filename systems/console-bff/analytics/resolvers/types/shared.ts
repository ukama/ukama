/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Shared analytics types reused across the business / customer / network /
 * collector backends.
 *
 * The analytics gateway renders protobuf responses with protojson
 * (EmitUnpopulated), so the REST JSON keys are already lowerCamelCase and
 * timestamps are RFC3339 strings — these GraphQL fields map 1:1 to the wire
 * shape, no snake_case remapping needed.
 */
import { Field, Float, InputType, Int, ObjectType } from "type-graphql";

@ObjectType()
export class KpiDto {
  @Field()
  key: string;

  @Field(() => Float)
  value: number;

  @Field({ nullable: true })
  formatted?: string;

  @Field(() => Float, { nullable: true })
  delta?: number;

  @Field({ nullable: true })
  deltaPeriod?: string;

  @Field({ nullable: true })
  stale?: boolean;

  @Field({ nullable: true })
  asOf?: string;
}

@ObjectType()
export class MetaDto {
  @Field(() => Int)
  count: number;

  @Field(() => Int)
  page: number;

  @Field(() => Int)
  size: number;

  @Field(() => Int)
  pages: number;
}

@ObjectType()
export class PointDto {
  @Field({ nullable: true })
  time?: string;

  @Field(() => Float)
  value: number;
}

@ObjectType()
export class TimeSeriesDto {
  @Field()
  key: string;

  @Field(() => [PointDto])
  points: PointDto[];
}

@ObjectType()
export class NamedValueDto {
  @Field({ nullable: true })
  name?: string;

  @Field({ nullable: true })
  id?: string;

  @Field(() => Float)
  value: number;
}

@ObjectType()
export class ActivityItemDto {
  @Field({ nullable: true })
  routingKey?: string;

  @Field({ nullable: true })
  description?: string;

  @Field({ nullable: true })
  occurredAt?: string;
}

/**
 * Generic windowed query input shared by most analytics queries. Every field
 * is optional; the gateway applies sensible defaults and ignores params a
 * given endpoint does not use (e.g. billing ignores networkId).
 */
@InputType()
export class AnalyticsWindowInput {
  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  siteId?: string;

  @Field({ nullable: true })
  nodeId?: string;

  @Field({ nullable: true })
  status?: string;

  @Field({ nullable: true })
  severity?: string;

  @Field({ nullable: true })
  state?: string;

  @Field({ nullable: true })
  metric?: string;

  @Field({ nullable: true })
  query?: string;

  @Field({ nullable: true })
  period?: string;

  @Field({ nullable: true })
  from?: string;

  @Field({ nullable: true })
  to?: string;

  @Field({ nullable: true })
  timezone?: string;

  @Field(() => Int, { nullable: true })
  page?: number;

  @Field(() => Int, { nullable: true })
  pageSize?: number;
}
