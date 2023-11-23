/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Ctx, Query, Resolver } from "type-graphql";
import { Arg, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { InvitationDto } from "./types";

@Resolver()
export class GetInvitationResolver {
  @Query(() => InvitationDto)
  @UseMiddleware(Authentication)
  async getInvitation(
    @Arg("id") id: string,
    @Ctx() ctx: Context
  ): Promise<InvitationDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getInvitation(id);
  }
}
