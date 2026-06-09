/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { InvitationsResDto } from "./types";

@Resolver()
export class GetInvitationsResolver {
  @Query(() => InvitationsResDto)
  async getInvitations(@Ctx() ctx: AppContext): Promise<InvitationsResDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("invitation");
    return await dataSources.invitation.getInvitations(baseURL);
  }
}
