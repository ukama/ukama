/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { SubscriberDto } from "./types";

@Resolver()
export class GetSubscriberResolver {
  @Query(() => SubscriberDto)
  async getSubscriber(
    @Arg("subscriberId") subscriberId: string,
    @Ctx() ctx: AppContext
  ): Promise<SubscriberDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("subscriber");
    return await dataSources.subscriber.getSubscriber(baseURL, subscriberId);
  }
}
