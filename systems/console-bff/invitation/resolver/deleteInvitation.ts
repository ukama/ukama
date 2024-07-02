/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { DeleteInvitationResDto } from "./types";

@Resolver()
export class DeleteInvitationResolver {
  @Mutation(() => DeleteInvitationResDto)
  async deleteInvitation(
    @Arg("id") id: string,
    @Ctx() ctx: Context
  ): Promise<DeleteInvitationResDto> {
    const { dataSources, baseURL } = ctx;
    return await dataSources.dataSource.deleteInvitation(baseURL, id);
  }
}
