/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import "reflect-metadata";
import { ArgsType, Field, InputType, ObjectType } from "type-graphql";

import {
  GRAPHS_TYPE,
  NOTIFICATION_SCOPE,
  NOTIFICATION_TYPE,
  ROLE_TYPE,
} from "../../common/enums";

@ObjectType()
export class LatestMetricRes {
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
@ObjectType()
export class MetricRes {
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

  @Field(() => [[Number, Number]])
  values: [number, number][];
}
@ArgsType()
@InputType()
export class GetLatestMetricInput {
  @Field()
  nodeId: string;

  @Field()
  type: string;
}

@ArgsType()
@InputType()
export class GetMetricRangeInput {
  @Field()
  nodeId: string;

  @Field()
  orgId?: string;

  @Field()
  type: string;

  @Field()
  userId?: string;

  @Field({ nullable: true })
  from: number;

  @Field({ nullable: true })
  to?: number;

  @Field({ nullable: true })
  step?: number;

  @Field({ nullable: true })
  withSubscription?: boolean;
}

@ArgsType()
@InputType()
export class GetMetricByTabInput {
  @Field()
  nodeId: string;

  @Field()
  orgId?: string;

  @Field(() => GRAPHS_TYPE)
  type: GRAPHS_TYPE;

  @Field()
  userId?: string;

  @Field({ nullable: true })
  from: number;

  @Field({ nullable: true })
  to?: number;

  @Field({ nullable: true })
  step?: number;

  @Field({ nullable: false })
  withSubscription?: boolean;
}

@ObjectType()
export class MetricsRes {
  @Field(() => [MetricRes])
  metrics: MetricRes[];
}

@ArgsType()
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

@ArgsType()
@InputType()
export class SubMetricByTabInput {
  @Field()
  nodeId: string;

  @Field()
  orgId: string;

  @Field(() => GRAPHS_TYPE)
  type: GRAPHS_TYPE;

  @Field()
  userId: string;

  @Field()
  from: number;
}

@ObjectType()
export class StatsMetric {
  @Field()
  activeSubscriber: number;

  @Field()
  averageSignalStrength: number;

  @Field()
  averageThroughput: number;
}

@ObjectType()
export class NotificationRes {
  @Field()
  id: string;

  @Field()
  title: string;

  @Field()
  description: string;

  @Field()
  orgId: string;

  @Field()
  networkId: string;

  @Field()
  subscriberId: string;

  @Field()
  userId: string;

  @Field()
  isRead: boolean;

  @Field(() => ROLE_TYPE)
  role: ROLE_TYPE;

  @Field(() => NOTIFICATION_TYPE)
  type: NOTIFICATION_TYPE;

  @Field(() => NOTIFICATION_SCOPE)
  scope: NOTIFICATION_SCOPE;
}

@ObjectType()
export class NotificationsRes {
  @Field(() => [NotificationRes])
  notifications: NotificationRes[];
}

@ArgsType()
@InputType()
export class GetNotificationsInput {
  @Field()
  orgId: string;

  @Field()
  networkId: string;

  @Field()
  siteId: string;

  @Field()
  nodeId: string;

  @Field()
  userId: string;

  @Field()
  subscriberId: string;

  @Field(() => ROLE_TYPE)
  forRole: ROLE_TYPE;
}

@ArgsType()
@InputType()
export class GetOrgNotificationsInput {
  @Field()
  orgId: string;

  @Field()
  userId: string;

  @Field(() => ROLE_TYPE)
  role: ROLE_TYPE;
}
@ArgsType()
@InputType()
export class GetNetworkNotificationsInput {
  @Field()
  orgId: string;

  @Field()
  userId: string;

  @Field()
  networkId: string;

  @Field(() => ROLE_TYPE)
  role: ROLE_TYPE;
}
@ArgsType()
@InputType()
export class GetSiteNotificationsInput {
  @Field()
  orgId: string;

  @Field()
  userId: string;

  @Field()
  networkId: string;

  @Field()
  siteId: string;

  @Field(() => ROLE_TYPE)
  role: ROLE_TYPE;
}
@ArgsType()
@InputType()
export class GetNodeNotificationsInput {
  @Field()
  orgId: string;

  @Field()
  userId: string;

  @Field()
  networkId: string;

  @Field()
  siteId: string;

  @Field()
  nodeId: string;

  @Field(() => ROLE_TYPE)
  role: ROLE_TYPE;
}
@ArgsType()
@InputType()
export class GetSubscriberNotificationsInput {
  @Field()
  orgId: string;

  @Field()
  userId: string;

  @Field()
  networkId: string;

  @Field()
  subscriberId: string;
}

@ArgsType()
@InputType()
export class GetUserNotificationsInput {
  @Field()
  userId: string;
}
