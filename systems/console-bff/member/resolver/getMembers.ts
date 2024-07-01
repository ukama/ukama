/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import axios from "axios";
import { Ctx, Query, Resolver } from "type-graphql";

import { NUCLEUS_API_GW } from "../../common/configs";
import { Context } from "../context";
import { dtoToUserResDto } from "../datasource/mapper";
import { MemberDto, MembersResDto } from "./types";

@Resolver()
export class GetMembersResolver {
  @Query(() => MembersResDto)
  async getMembers(@Ctx() ctx: Context): Promise<MembersResDto> {
    const { dataSources, baseURL } = ctx;
    const members: MemberDto[] = [];
    const res = await dataSources.dataSource.getMembers(baseURL);
    for (const member of res.members) {
      const user = await axios
        .get(`${NUCLEUS_API_GW}/v1/users/${member.userId}`)
        .then(res => dtoToUserResDto(res.data));

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
