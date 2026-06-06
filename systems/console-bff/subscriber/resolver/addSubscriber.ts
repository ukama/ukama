/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { SubscriberDto, SubscriberInputDto } from "./types";

@Resolver()
export class AddSubscriberResolver {
  @Mutation(() => SubscriberDto)
  async addSubscriber(
    @Arg("data") data: SubscriberInputDto,
    @Ctx() ctx: AppContext
  ): Promise<SubscriberDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("subscriber");
    return await dataSources.subscriber.addSubscriber(baseURL, data);
  }
}
