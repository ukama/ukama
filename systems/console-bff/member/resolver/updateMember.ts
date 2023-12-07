/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { UpdateMemberInputDto } from "./types";

@Resolver()
export class UpdateMemberResolver {
  @Mutation(() => CBooleanResponse)
  async updateMember(
    @Arg("memberId") memberId: string,
    @Arg("data") data: UpdateMemberInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.updateMember(memberId, data);
  }
}
