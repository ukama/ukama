import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { MemberObj, OrgMembersResDto } from "./types";

@Resolver()
export class GetOrgMembersResolver {
  @Query(() => OrgMembersResDto)
  async getOrgMembers(@Ctx() ctx: Context): Promise<OrgMembersResDto> {
    const { dataSources, headers } = ctx;
    const res: MemberObj[] = [];
    const members: OrgMembersResDto =
      await dataSources.dataSource.getOrgMembers(headers.orgName);

    if (members.members.length === 0) return members;
    else {
      for (const member of members.members) {
        // const user = await dataSources.dataSoureceUser.getUser(member.uuid);
        res.push({
          ...member,
          // user: {
          //   name: user.name,
          //   uuid: user.uuid,
          //   email: user.email,
          //   phone: user.phone,
          //   isDeactivated: user.isDeactivated,
          //   registeredSince: user.registeredSince,
          // },
        });
      }
      return { members: res, org: members.org };
    }
  }
}
