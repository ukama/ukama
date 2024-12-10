/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, ObjectType } from "type-graphql";

import { NOTIFICATION_SCOPE, NOTIFICATION_TYPE } from "../../common/enums";

@ObjectType()
export class NotificationsAPIDto {
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
  created_at: string;
}

@ObjectType()
export class NotificationAPIDto {
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
  created_at: string;

  @Field()
  org_id: string;

  @Field()
  network_id: string;

  @Field()
  subscriber_id: string;

  @Field()
  user_id: string;

  @Field()
  resource_id: string;
}

@ObjectType()
export class NotificationAPIRes {
  @Field(() => NotificationAPIDto)
  notification: NotificationAPIDto;
}
@ObjectType()
export class NotificationsAPIRes {
  @Field(() => [NotificationsAPIDto])
  notifications: NotificationsAPIDto[];
}

@ObjectType()
export class NotificationResDto {
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
  orgId: string;

  @Field()
  userId: string;

  @Field()
  subscriberId: string;

  @Field()
  networkId: string;

  @Field()
  createdAt: string;

  @Field()
  resourceId: string;
}

@ObjectType()
export class NotificationsDto {
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
  isRead: boolean;

  @Field()
  createdAt: string;
}

@ObjectType()
export class UpdateNotificationResDto {
  @Field()
  id: string;
}

@ObjectType()
export class NotificationsResDto {
  @Field(() => [NotificationsDto])
  notifications: NotificationsDto[];
}
