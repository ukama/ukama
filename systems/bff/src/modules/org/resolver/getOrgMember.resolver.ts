import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { OrgService } from "../service";
import { MemberObj } from "../types";

@Service()
@Resolver()
export class GetOrgMemberResolver {
    constructor(private readonly orgService: OrgService) {}

    @Query(() => MemberObj)
    @UseMiddleware(Authentication)
    async getOrgMember(@Ctx() ctx: Context): Promise<MemberObj> {
        return this.orgService.getOrgMember(parseHeaders(ctx));
    }
}
