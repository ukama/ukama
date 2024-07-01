/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddMemberInputDto, MemberDto } from "./types";

@Resolver()
export class AddMemberResolver {
  @Mutation(() => MemberDto)
  async addMember(
    @Arg("data") data: AddMemberInputDto,
    @Ctx() ctx: Context
  ): Promise<MemberDto> {
    const { dataSources, baseURL } = ctx;
    return dataSources.dataSource.addMember(baseURL, data);
  }
}
