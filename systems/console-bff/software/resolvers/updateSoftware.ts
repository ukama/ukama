/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { StringResponse, UpdateSoftwareInputDto } from "./types";

@Resolver()
export class UpdateSoftwareResolver {
  @Mutation(() => StringResponse)
  async updateSoftware(
    @Arg("data") data: UpdateSoftwareInputDto,
    @Ctx() ctx: AppContext
  ): Promise<StringResponse> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("software");
    return dataSources.software.updateSoftware(baseURL, data);
  }
}
