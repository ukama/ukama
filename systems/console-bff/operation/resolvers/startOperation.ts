/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { StartOperationInputDto, StartOperationResponseDto } from "./types";

@Resolver()
export class StartOperationResolver {
  @Mutation(() => StartOperationResponseDto)
  async startOperation(
    @Arg("data") data: StartOperationInputDto,
    @Ctx() ctx: AppContext
  ): Promise<StartOperationResponseDto> {
    const baseURL = await ctx.urls.url("operation");
    return ctx.dataSources.operation.startOperation(baseURL, data);
  }
}
