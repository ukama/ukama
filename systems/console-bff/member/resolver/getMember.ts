/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { MemberDto } from "./types";

@Resolver()
export class GetMemberResolver {
  @Query(() => MemberDto)
  async getMember(
    @Arg("id") id: string,
    @Ctx() ctx: Context
  ): Promise<MemberDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getMember(id);
  }
}
