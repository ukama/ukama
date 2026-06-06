/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { AddMemberInputDto, MemberDto } from "./types";

@Resolver()
export class AddMemberResolver {
  @Mutation(() => MemberDto)
  async addMember(
    @Arg("data") data: AddMemberInputDto,
    @Ctx() ctx: AppContext
  ): Promise<MemberDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("member");
    return dataSources.member.addMember(baseURL, data);
  }
}
