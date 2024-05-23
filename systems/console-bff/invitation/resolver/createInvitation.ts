/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { CreateInvitationInputDto, InvitationDto } from "./types";

@Resolver()
export class CreateInvitationResolver {
  @Mutation(() => InvitationDto)
  async createInvitation(
    @Arg("data") data: CreateInvitationInputDto,
    @Ctx() ctx: Context
  ): Promise<InvitationDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.sendInvitation(data);
  }
}
