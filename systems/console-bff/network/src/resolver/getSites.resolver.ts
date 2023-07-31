import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SitesResDto } from "../types";

@Resolver()
export class GetSitesResolver {

    @Query(() => SitesResDto)
    @UseMiddleware(Authentication)
    async getSites(
        @Arg("networkId") networkId: string,
        @Ctx() ctx: Context
    ): Promise<SitesResDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getSites(networkId, parseHeaders(ctx));
    }
}
