/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import axios from "axios";
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { NUCLEUS_API_GW } from "../../common/configs";
import { Context } from "../context";
import { dtoToUserResDto } from "../datasource/mapper";
import { MemberDto } from "./types";

@Resolver()
export class GetMemberByUserIdResolver {
  @Query(() => MemberDto)
  async getMemberByUserId(
    @Arg("userId") id: string,
    @Ctx() ctx: Context
  ): Promise<MemberDto> {
    const { dataSources, baseURL } = ctx;
    const member = await dataSources.dataSource.getMemberByUserId(baseURL, id);

    const user = await axios
      .get(`${NUCLEUS_API_GW}/v1/users/${member.userId}`)
      .then(res => dtoToUserResDto(res.data));

    return {
      role: member.role,
      userId: member.userId,
      name: user.name ?? "",
      email: user.email ?? "",
      memberId: member.memberId,
      memberSince: member.memberSince,
      isDeactivated: member.isDeactivated,
    };
  }
}
