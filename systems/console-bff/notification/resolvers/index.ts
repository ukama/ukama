/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { GetNotificationResolver } from "./getNotification";
import { GetNotificationsResolver } from "./getNotifications";
import { UpdateNotificationResolver } from "./updateNotification";

const resolvers: NonEmptyArray<any> = [
  GetNotificationResolver,
  GetNotificationsResolver,
  UpdateNotificationResolver,
];

export default resolvers;
