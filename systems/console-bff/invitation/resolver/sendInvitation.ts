/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { SendInvitationInputDto, SendInvitationResDto } from "./types";

@Resolver()
export class SendInvitationResolver {
  @Mutation(() => SendInvitationResDto)
  @UseMiddleware(Authentication)
  async sendInvitation(
    @Arg("data") data: SendInvitationInputDto,
    @Ctx() ctx: Context
  ): Promise<SendInvitationResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.sendInvitation(data);
  }
}
