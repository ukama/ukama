/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { NotificationsResDto } from "./types";

@Resolver()
export class GetNotificationsResolver {
  @Query(() => NotificationsResDto)
  async getNotifications(@Ctx() ctx: Context): Promise<NotificationsResDto> {
    const { dataSources, baseURL, headers } = ctx;
    return dataSources.dataSource.getNotifications(
      baseURL,
      headers.orgId,
      headers.userId
    );
  }
}
