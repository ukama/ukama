import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { MemberObj, OrgMembersResDto } from "./types";

@Resolver()
export class GetOrgMembersResolver {
  @Query(() => OrgMembersResDto)
  @UseMiddleware(Authentication)
  async getOrgMembers(@Ctx() ctx: Context): Promise<OrgMembersResDto> {
    const { dataSources } = ctx;

    const res: MemberObj[] = [];
    const members = await dataSources.dataSource.getOrgMembers(parseHeaders());
    for (const member of members.members) {
      const user = await getUser(member.uuid, parseHeaders());
      res.push({
        ...member,
        user: {
          name: user.name,
          uuid: user.uuid,
          email: user.email,
          phone: user.phone,
          isDeactivated: user.isDeactivated,
          registeredSince: user.registeredSince,
        },
      });
    }
    return { members: res, org: members.org };
  }
}
