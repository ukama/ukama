/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { NotificationsResDto } from "./types";

@Resolver()
export class GetNotificationsResolver {
  @Query(() => NotificationsResDto)
  async getNotifications(@Ctx() ctx: AppContext): Promise<NotificationsResDto> {
    const { dataSources, headers } = ctx;
    const baseURL = await ctx.urls.url("notification");
    return dataSources.notification.getNotifications(
      baseURL,
      headers.orgId,
      headers.userId
    );
  }
}
