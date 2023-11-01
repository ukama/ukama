/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { SubscriberDto, SubscriberInputDto } from "./types";

@Resolver()
export class AddSubscriberResolver {
  @Mutation(() => SubscriberDto)
  async addSubscriber(
    @Arg("data") data: SubscriberInputDto,
    @Ctx() ctx: Context
  ): Promise<SubscriberDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.addSubscriber(data);
  }
}
