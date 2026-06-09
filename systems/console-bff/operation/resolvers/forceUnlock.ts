/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { ForceUnlockInputDto, OperationDto } from "./types";

@Resolver()
export class ForceUnlockResolver {
  @Mutation(() => OperationDto, { nullable: true })
  async forceUnlock(
    @Arg("data") data: ForceUnlockInputDto,
    @Ctx() ctx: AppContext
  ): Promise<OperationDto | undefined> {
    // The gateway role-gates force-unlock to org owner/admin using the caller's
    // user id; take it from the verified token claims, never from client input.
    const baseURL = await ctx.urls.url("operation");
    return ctx.dataSources.operation.forceUnlock(
      baseURL,
      data,
      ctx.headers.userId ?? ""
    );
  }
}
