/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import { INVITATION_STATUS } from "../../common/enums";
import { Context } from "../context";
import { InvitationDto, InvitationsResDto } from "./types";

@Resolver()
export class GetInVitationsByOrgResolver {
  @Query(() => InvitationsResDto)
  async getInvitationsByOrg(@Ctx() ctx: Context): Promise<InvitationsResDto> {
    const { dataSources } = ctx;
    const res = await dataSources.dataSource.getInvitationsByOrg();
    const Invitations: InvitationDto[] = [];
    for (const invitation of res.invitations) {
      if (invitation.status !== INVITATION_STATUS.ACCEPTED) {
        Invitations.push(invitation);
      }
    }
    return {
      invitations: Invitations,
    };
  }
}
