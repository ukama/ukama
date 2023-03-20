import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { OrgService } from "../service";
import { OrgMembersResDto } from "../types";

@Service()
@Resolver()
export class GetOrgMembersResolver {
    constructor(private readonly orgService: OrgService) {}

    @Query(() => OrgMembersResDto)
    @UseMiddleware(Authentication)
    async getOrgMembers(@Ctx() ctx: Context): Promise<OrgMembersResDto> {
        return this.orgService.getOrgMembers(parseCookie(ctx));
    }
}
