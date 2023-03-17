import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { OrgService } from "../service";
import { OrgsResDto } from "../types";

@Service()
@Resolver()
export class GetOrgsResolver {
    constructor(private readonly orgService: OrgService) {}

    @Query(() => OrgsResDto)
    @UseMiddleware(Authentication)
    async getOrgs(@Ctx() ctx: Context): Promise<OrgsResDto> {
        return this.orgService.getOrgs(parseCookie(ctx));
    }
}
