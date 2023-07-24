import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../../../user/src/service";
import { MemberObj, OrgMembersResDto } from "../types";

@Resolver()
export class GetOrgMembersResolver {

    @Query(() => OrgMembersResDto)
    @UseMiddleware(Authentication)
    async getOrgMembers(@Ctx() ctx: Context): Promise<OrgMembersResDto> {
        const { dataSources } = ctx;

        const res: MemberObj[] = [];
        const members = await dataSources.dataSource.getOrgMembers(parseHeaders(ctx));
        for (const member of members.members) {
            const userService: UserService = new UserService();
            const user = await userService.getUser(
                member.uuid,
                parseHeaders(ctx)
            );
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
