/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { UpdateInvitationInputDto, UpdateInvitationResDto } from "./types";

@Resolver()
export class UpdateInvitationResolver {
  @Mutation(() => UpdateInvitationResDto)
  async updateInvitation(
    @Arg("data") data: UpdateInvitationInputDto,
    @Ctx() ctx: AppContext
  ): Promise<UpdateInvitationResDto> {
    const { dataSources } = ctx;
    return await dataSources.invitation.updateInvitation(data);
  }
}
