/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { ArgsType, Field, InputType, ObjectType } from "type-graphql";

import {
  GRAPHS_TYPE,
  NOTIFICATION_SCOPE,
  NOTIFICATION_TYPE,
  ROLE_TYPE,
} from "../../common/enums";

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
@InputType()
export class GetLatestMetricInput {
  @Field()
  nodeId: string;

  @Field()
  type: string;
}

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

@InputType()
export class GetMetricByTabInput {
  @Field()
  nodeId: string;

  @Field()
  orgId?: string;

  @Field()
  orgName: string;

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
export class NotificationsAPIResDto {
  @Field()
  id: string;

  @Field()
  title: string;

  @Field()
  description: string;

  @Field(() => NOTIFICATION_TYPE)
  type: NOTIFICATION_TYPE;

  @Field(() => NOTIFICATION_SCOPE)
  scope: NOTIFICATION_SCOPE;

  @Field()
  is_read: boolean;

  @Field()
  is_actionable: boolean;

  @Field()
  resource_id: string;

  @Field()
  created_at: string;
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
  createdAt: string;

  @Field(() => NOTIFICATION_TYPE)
  type: NOTIFICATION_TYPE;

  @Field(() => NOTIFICATION_SCOPE)
  scope: NOTIFICATION_SCOPE;

  @Field()
  isRead: boolean;

  @Field()
  isActionable: boolean;

  @Field()
  resourceId: string;
}

@ObjectType()
export class NotificationsRes {
  @Field(() => [NotificationsResDto])
  notifications: NotificationsResDto[];
}

@ArgsType()
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
