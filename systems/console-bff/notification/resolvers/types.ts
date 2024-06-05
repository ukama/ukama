/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, ObjectType } from "type-graphql";

import {
  NOTIFICATION_SCOPE,
  NOTIFICATION_TYPE,
  ROLE_TYPE,
} from "../../common/enums";

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

  @Field(() => ROLE_TYPE)
  for_role: ROLE_TYPE;

  @Field()
  created_at: string;
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

  @Field(() => NOTIFICATION_TYPE)
  type: NOTIFICATION_TYPE;

  @Field(() => NOTIFICATION_SCOPE)
  scope: NOTIFICATION_SCOPE;

  @Field(() => ROLE_TYPE)
  forRole: ROLE_TYPE;

  @Field()
  createdAt: string;
}

@ObjectType()
export class UpdateNotificationResDto {
  @Field()
  id: string;
}
