/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class NetworkStats {
  @Field()
  activeSubscriber: number;

  @Field()
  averageSignalStrength: number;

  @Field()
  averageThroughput: number;
}
@InputType()
export class GetNodeLatestMetricInput {
  @Field()
  nodeId: string;

  @Field()
  type: string;
}

@InputType()
export class GetSiteLatestMetricInput {
  @Field()
  siteId: string;

  @Field()
  type: string;
}

@ObjectType()
export class SiteLatestMetric {
  @Field()
  success: boolean;

  @Field()
  msg: string;

  @Field()
  orgId: string;

  @Field()
  siteId: string;

  @Field()
  type: string;

  @Field(() => [Number, Number])
  value: [number, number];
}

@ObjectType()
export class NodeLatestMetric {
  @Field()
  success: boolean;

  @Field()
  msg: string;

  @Field()
  orgId: string;

  @Field()
  nodeId: string;

  @Field()
  type: string;

  @Field(() => [Number, Number])
  value: [number, number];
}
