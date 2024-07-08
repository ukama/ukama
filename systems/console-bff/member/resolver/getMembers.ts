/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import UserApi from "../datasource/user_api";
import { MemberDto, MembersResDto } from "./types";

@Resolver()
export class GetMembersResolver {
  @Query(() => MembersResDto)
  async getMembers(@Ctx() ctx: Context): Promise<MembersResDto> {
    const { dataSources, baseURL } = ctx;
    const members: MemberDto[] = [];
    const res = await dataSources.dataSource.getMembers(baseURL);
    const userAPI = new UserApi();
    for (const member of res.members) {
      const user = await userAPI.getUser(member.userId);

      members.push({
        role: member.role,
        userId: member.userId,
        name: user.name ?? "",
        email: user.email ?? "",
        memberId: member.memberId,
        memberSince: member.memberSince,
        isDeactivated: member.isDeactivated,
      });
    }

    return {
      members: members,
    };
  }
}
