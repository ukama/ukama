/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, ObjectType } from "type-graphql";

@ObjectType()
export class MetricAggregated {
  @Field()
  computed_at: string;

  @Field()
  value: number;

  @Field()
  min: number;

  @Field()
  max: number;

  @Field()
  p95: number;

  @Field()
  mean: number;

  @Field()
  median: number;

  @Field()
  sample_count: number;

  @Field()
  aggregation: string;

  @Field()
  noise_estimate: number;
}

@ObjectType()
export class MetricTrend {
  @Field()
  type: string;

  @Field()
  value: number;
}

@ObjectType()
export class MetricConfidence {
  @Field()
  value: number;
}

@ObjectType()
export class MetricProjection {
  @Field()
  type: string;

  @Field()
  eta_sec: number;
}

@ObjectType()
export class MetricAnalysis {
  @Field(() => MetricAggregated)
  aggregated: MetricAggregated;

  @Field(() => MetricTrend)
  trend: MetricTrend;

  @Field(() => MetricConfidence)
  confidence: MetricConfidence;

  @Field(() => MetricProjection)
  projection: MetricProjection;

  @Field()
  state: string;
}

@ObjectType()
export class MetricDomain {
  @Field()
  rule_id: string;

  @Field()
  severity: string;

  @Field()
  headline: string;

  @Field()
  root_cause: string;

  @Field()
  service_impact: string;

  @Field()
  rule_confidence: number;

  @Field()
  evaluated_at: string;

  @Field()
  computed_at: string;
}
