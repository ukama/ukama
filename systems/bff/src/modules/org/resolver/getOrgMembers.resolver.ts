import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../../user/service";
import { OrgService } from "../service";
import { MemberObj, OrgMembersResDto } from "../types";

@Service()
@Resolver()
export class GetOrgMembersResolver {
    constructor(private readonly orgService: OrgService) {}

    @Query(() => OrgMembersResDto)
    @UseMiddleware(Authentication)
    async getOrgMembers(@Ctx() ctx: Context): Promise<OrgMembersResDto> {
        const res: MemberObj[] = [];
        const members = await this.orgService.getOrgMembers(parseHeaders(ctx));
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
