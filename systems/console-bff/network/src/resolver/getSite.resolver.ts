import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SiteDto } from "./../types";

@Resolver()
export class GetSiteResolver {

    @Query(() => SiteDto)
    @UseMiddleware(Authentication)
    async getSite(
        @Arg("networkId") networkId: string,
        @Arg("siteId") siteId: string,
        @Ctx() ctx: Context
    ): Promise<SiteDto> {
        const { dataSources } = ctx;

        return dataSources.dataSource.getSite(
            networkId,
            siteId,
            parseHeaders(ctx)
        );
    }
}
