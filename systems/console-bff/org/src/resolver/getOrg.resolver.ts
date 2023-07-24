import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { OrgDto } from "../types";

@Resolver()
export class GetOrgResolver {

    @Query(() => OrgDto)
    @UseMiddleware(Authentication)
    async getOrg(
        @Arg("orgName") orgName: string,
        @Ctx() ctx: Context
    ): Promise<OrgDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getOrg(orgName, parseHeaders(ctx));
    }
}
