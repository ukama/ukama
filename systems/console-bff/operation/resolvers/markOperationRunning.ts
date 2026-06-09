/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { MarkOperationRunningInputDto, OperationDto } from "./types";

@Resolver()
export class MarkOperationRunningResolver {
  @Mutation(() => OperationDto, { nullable: true })
  async markOperationRunning(
    @Arg("data") data: MarkOperationRunningInputDto,
    @Ctx() ctx: AppContext
  ): Promise<OperationDto | undefined> {
    const baseURL = await ctx.urls.url("operation");
    return ctx.dataSources.operation.markOperationRunning(baseURL, data);
  }
}
