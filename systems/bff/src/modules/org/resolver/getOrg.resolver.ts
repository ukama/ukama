import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { OrgService } from "../service";
import { OrgDto } from "../types";

@Service()
@Resolver()
export class GetOrgResolver {
    constructor(private readonly orgService: OrgService) {}

    @Query(() => OrgDto)
    @UseMiddleware(Authentication)
    async getOrg(
        @Arg("orgName") orgName: string,
        @Ctx() ctx: Context,
    ): Promise<OrgDto> {
        return this.orgService.getOrg(orgName, parseCookie(ctx));
    }
}
