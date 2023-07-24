import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { OrgsResDto } from "../types";

@Resolver()
export class GetOrgsResolver {

    @Query(() => OrgsResDto)
    @UseMiddleware(Authentication)
    async getOrgs(@Ctx() ctx: Context): Promise<OrgsResDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getOrgs(parseHeaders(ctx));
    }
}
