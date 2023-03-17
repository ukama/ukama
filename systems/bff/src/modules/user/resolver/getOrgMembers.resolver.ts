import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { OrgMemberResDto } from "../types";

@Service()
@Resolver()
export class GetOrgMembersResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => OrgMemberResDto)
    @UseMiddleware(Authentication)
    async getOrgMembers(@Ctx() ctx: Context): Promise<OrgMemberResDto> {
        return this.userService.getUsersByOrg(parseCookie(ctx));
    }
}
