/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { INVITATION_STATUS } from "../../common/enums";
import { Context } from "../context";
import { InvitationDto, InvitationsResDto } from "./types";

@Resolver()
export class GetInVitationsByEmailResolver {
  @Query(() => InvitationsResDto)
  async getInvitationsByEmail(
    @Arg("email") email: string,
    @Ctx() ctx: Context
  ): Promise<InvitationsResDto> {
    const { dataSources } = ctx;
    const res = await dataSources.dataSource.getAllInvitationsByEmail(email);
    const Invitations: InvitationDto[] = [];
    for (const invitation of res.invitations) {
      if (invitation.status !== INVITATION_STATUS.INVITE_ACCEPTED) {
        Invitations.push(invitation);
      }
    }
    return {
      invitations: Invitations,
    };
  }
}
