/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetOperationInputDto, OperationDto } from "./types";

@Resolver()
export class GetOperationResolver {
  @Query(() => OperationDto, { nullable: true })
  async getOperation(
    @Arg("data") data: GetOperationInputDto,
    @Ctx() ctx: AppContext
  ): Promise<OperationDto | undefined> {
    const baseURL = await ctx.urls.url("operation");
    return ctx.dataSources.operation.getOperation(baseURL, data.id);
  }
}
