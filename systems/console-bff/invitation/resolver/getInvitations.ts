/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { InvitationDto } from "./types";

@Resolver()
export class GetInvitationsResolver {
  @Query(() => InvitationDto)
  async getInvitations(
    @Arg("email") email: string,
    @Ctx() ctx: Context
  ): Promise<InvitationDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getInvitationsByEmail(email);
  }
}
