/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, ObjectType } from "type-graphql";

@ObjectType()
export class NotificationAPIDto {
  @Field()
  id: string;

  @Field()
  title: string;

  @Field()
  description: string;

  @Field()
  type: string;

  @Field()
  scope: string;

  @Field()
  org_id: string;

  @Field()
  network_id: string;

  @Field()
  subscriber_id: string;

  @Field()
  user_id: string;

  @Field()
  for_role: string;
}

@ObjectType()
export class NotificationAPIRes {
  @Field(() => NotificationAPIDto)
  notification: NotificationAPIDto;
}

@ObjectType()
export class NotificationResDto {
  @Field()
  id: string;

  @Field()
  title: string;

  @Field()
  description: string;

  @Field()
  type: string;

  @Field()
  scope: string;

  @Field()
  orgId: string;

  @Field()
  networkId: string;

  @Field()
  subscriberId: string;

  @Field()
  userId: string;

  @Field()
  forRole: string;
}

@ObjectType()
export class UpdateNotificationResDto {
  @Field()
  id: string;
}
