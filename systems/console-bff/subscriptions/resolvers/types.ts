/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

import { METRICS_INTERVAL } from "../../common/configs";
import {
  GRAPHS_TYPE,
  NOTIFICATION_SCOPE,
  NOTIFICATION_TYPE,
  ROLE_TYPE,
  STATS_TYPE,
} from "../../common/enums";

@ObjectType()
export class MetricRes {
  @Field()
  success: boolean;

  @Field()
  msg: string;

  @Field({ nullable: true })
  nodeId?: string;

  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  packageId?: string;

  @Field({ nullable: true })
  dataPlanId?: string;

  @Field({ nullable: true })
  siteId?: string;

  @Field()
  type: string;

  @Field(() => [[Number, Number]])
  values: [number, number][];
}
@InputType()
export class GetLatestMetricInput {
  @Field()
  nodeId: string;

  @Field()
  type: string;

  @Field()
  siteId?: string;
}

@ObjectType()
export class LatestMetricRes {
  @Field()
  success: boolean;

  @Field()
  msg: string;

  @Field()
  nodeId: string;

  @Field()
  siteId?: string;

  @Field()
  type: string;

  @Field(() => [Number, Number])
  value: [number, number];
}

@ObjectType()
export class LatestMetricSubRes {
  @Field()
  success: boolean;

  @Field()
  msg: string;

  @Field()
  nodeId: string;

  @Field()
  siteId?: string;

  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  packageId?: string;

  @Field({ nullable: true })
  dataPlanId?: string;

  @Field()
  type: string;

  @Field(() => [Number, Number])
  value: [number, number];
}

@InputType()
export class GetMetricRangeInput {
  @Field()
  nodeId?: string;

  @Field()
  orgId?: string;

  @Field()
  siteId?: string;

  @Field()
  type: string;

  @Field()
  userId?: string;

  @Field({ nullable: true })
  from: number;

  @Field({ nullable: true })
  to?: number;

  @Field({ nullable: true, defaultValue: METRICS_INTERVAL })
  step?: number;

  @Field({ nullable: true })
  withSubscription?: boolean;
}

@InputType()
export class GetMetricsStatInput {
  @Field({ nullable: true })
  nodeId?: string;

  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  siteId?: string;

  @Field()
  orgName: string;

  @Field({ nullable: true })
  userId?: string;

  @Field(() => STATS_TYPE)
  type: STATS_TYPE;

  @Field({ nullable: true, defaultValue: "avg" })
  operation?: string;

  @Field()
  from: number;

  @Field({ nullable: true })
  to?: number;

  @Field({ defaultValue: 30 })
  step: number;

  @Field({ defaultValue: false })
  withSubscription: boolean;
}

@InputType()
export class GetMetricsSiteStatInput {
  @Field()
  orgName: string;

  @Field({ nullable: true })
  userId?: string;

  @Field(() => STATS_TYPE)
  type: STATS_TYPE;

  @Field(() => [String], { nullable: true })
  siteIds?: string[];

  @Field(() => [String], { nullable: true })
  nodeIds?: string[];

  @Field()
  from: number;

  @Field()
  to: number;

  @Field({ nullable: true, defaultValue: "avg" })
  operation?: string;

  @Field({ defaultValue: 30 })
  step: number;

  @Field({ defaultValue: false })
  withSubscription: boolean;
}

@InputType()
export class SubMetricsStatInput {
  @Field({ nullable: true })
  nodeId?: string;

  @Field({ nullable: true })
  networkId?: string;

  @Field()
  orgName: string;

  @Field(() => STATS_TYPE)
  type: STATS_TYPE;

  @Field()
  userId: string;

  @Field()
  from: number;
}

@InputType()
export class SubSiteMetricsStatInput {
  @Field(() => [String], { nullable: true })
  siteIds?: string[];

  @Field()
  orgName: string;

  @Field(() => STATS_TYPE)
  type: STATS_TYPE;

  @Field(() => [String], { nullable: true })
  nodeIds?: string[];

  @Field()
  userId: string;

  @Field()
  from: number;
}

@InputType()
export class GetMetricByTabInput {
  @Field({ nullable: true })
  nodeId?: string;

  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  siteId?: string;

  @Field()
  orgName: string;

  @Field()
  userId: string;

  @Field(() => GRAPHS_TYPE)
  type: GRAPHS_TYPE;

  @Field()
  from: number;

  @Field()
  to: number;

  @Field({ defaultValue: 1 })
  step: number;

  @Field({ defaultValue: false })
  withSubscription: boolean;
}

@ObjectType()
export class MetricsRes {
  @Field(() => [MetricRes])
  metrics: MetricRes[];
}

@ObjectType()
export class MetricStateRes {
  @Field()
  success: boolean;

  @Field()
  msg: string;

  @Field()
  nodeId?: string;

  @Field()
  type: string;

  @Field({ nullable: true })
  siteId?: string;

  @Field()
  value: number;

  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  packageId?: string;

  @Field({ nullable: true })
  dataPlanId?: string;
}

@InputType()
export class GetMetricBySiteInput {
  @Field()
  orgName: string;

  @Field()
  userId: string;

  @Field(() => GRAPHS_TYPE)
  type: GRAPHS_TYPE;

  @Field()
  from: number;

  @Field()
  siteId?: string;

  @Field()
  to: number;

  @Field({ defaultValue: 1 })
  step: number;

  @Field({ defaultValue: false })
  withSubscription: boolean;
}

@ObjectType()
export class MetricsStateRes {
  @Field(() => [MetricStateRes])
  metrics: MetricStateRes[];
}

@InputType()
export class SubMetricRangeInput {
  @Field()
  nodeId: string;

  @Field()
  orgId: string;

  @Field()
  type: string;

  @Field()
  userId: string;

  @Field()
  from: number;
}

@InputType()
export class SubMetricByTabInput {
  @Field()
  nodeId: string;

  @Field()
  orgName: string;

  @Field(() => GRAPHS_TYPE)
  type: GRAPHS_TYPE;

  @Field()
  userId: string;

  @Field()
  from: number;
}
@InputType()
export class SubSiteMetricByTabInput {
  @Field()
  siteId: string;

  @Field()
  orgName: string;

  @Field(() => GRAPHS_TYPE)
  type: GRAPHS_TYPE;

  @Field()
  userId: string;

  @Field()
  from: number;
}

@ObjectType()
export class NotificationsAPIResDto {
  @Field()
  id: string;

  @Field()
  title: string;

  @Field()
  description: string;

  @Field()
  resource_id: string;

  @Field()
  event_key: string;

  @Field(() => NOTIFICATION_TYPE)
  type: NOTIFICATION_TYPE;

  @Field(() => NOTIFICATION_SCOPE)
  scope: NOTIFICATION_SCOPE;

  @Field()
  is_read: boolean;

  @Field()
  created_at: string;
}

@ObjectType()
export class NotificationRedirect {
  @Field()
  title: string;

  @Field()
  action: string;
}
@ObjectType()
export class NotificationsAPIRes {
  @Field(() => [NotificationsAPIResDto])
  notifications: NotificationsAPIResDto[];
}

@ObjectType()
export class NotificationsResDto {
  @Field()
  id: string;

  @Field()
  title: string;

  @Field()
  description: string;

  @Field()
  eventKey: string;

  @Field()
  resourceId: string;

  @Field()
  createdAt: string;

  @Field(() => NOTIFICATION_TYPE)
  type: NOTIFICATION_TYPE;

  @Field(() => NOTIFICATION_SCOPE)
  scope: NOTIFICATION_SCOPE;

  @Field()
  isRead: boolean;

  @Field(() => NotificationRedirect)
  redirect?: NotificationRedirect;
}

@ObjectType()
export class NotificationsRes {
  @Field(() => [NotificationsResDto])
  notifications: NotificationsResDto[];
}

@InputType()
export class GetNotificationsInput {
  @Field({ nullable: false })
  orgName: string;

  @Field({ nullable: false })
  orgId: string;

  @Field({ nullable: false })
  userId: string;

  @Field({ nullable: false })
  startTimestamp: string;

  @Field()
  networkId: string;

  @Field()
  subscriberId: string;

  @Field(() => ROLE_TYPE)
  role: ROLE_TYPE;
}
