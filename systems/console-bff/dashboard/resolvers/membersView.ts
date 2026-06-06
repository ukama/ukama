/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Ctx, FieldResolver, Query, Resolver, Root } from "type-graphql";

import type { AppContext } from "../../server/context";
import { ServiceUrlResolver } from "../baseUrls";
import { runSection } from "../section";
import { MembersView, TeamMemberDto, TeamSection } from "./types";

type MembersRoot = MembersView & { _urls: ServiceUrlResolver };

/**
 * Members & admin composite (plan §3.2): org members and pending
 * invitations merged into one team list (status: Active | Deactivated |
 * Invited) so the screen renders a single table.
 */
@Resolver(() => MembersView)
export class MembersViewResolver {
  @Query(() => MembersView)
  membersView(@Ctx() ctx: AppContext): MembersView {
    return Object.assign(new MembersView(), {
      orgName: ctx.headers.orgName,
      _urls: new ServiceUrlResolver(ctx.headers.orgName),
    });
  }

  @FieldResolver(() => TeamSection)
  async team(
    @Root() root: MembersRoot,
    @Ctx() ctx: AppContext
  ): Promise<TeamSection> {
    const { value, error } = await runSection("team", async () => {
      const [memberUrl, invitationUrl] = await Promise.all([
        root._urls.url("member"),
        root._urls.url("invitation"),
      ]);
      const [members, invitations] = await Promise.all([
        ctx.dataSources.member.getMembers(memberUrl),
        ctx.dataSources.invitation.getInvitations(invitationUrl),
      ]);
      // The registry member list carries only userId+role; enrich each row
      // with the user's name/email (nucleus). allSettled so one failed
      // lookup leaves that row name-less rather than breaking the section.
      const users = await Promise.allSettled(
        members.members.map(m => ctx.dataSources.user.getUser(m.userId))
      );
      const rows: TeamMemberDto[] = members.members.map((member, i) => {
        const u = users[i];
        const user = u?.status === "fulfilled" ? u.value : undefined;
        return {
          id: member.memberId,
          name: user?.name ?? member.name,
          email: user?.email ?? member.email,
          role: member.role,
          status: member.isDeactivated ? "Deactivated" : "Active",
          memberSince: member.memberSince,
        };
      });
      for (const invitation of invitations.invitations) {
        rows.push({
          id: invitation.id,
          name: invitation.name,
          email: invitation.email,
          role: invitation.role,
          status: "Invited",
          inviteExpiresAt: invitation.expireAt,
        });
      }
      return rows;
    });
    return { rows: value, error };
  }
}
