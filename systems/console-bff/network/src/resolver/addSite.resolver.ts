import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddSiteInputDto, SiteDto } from "./../types";

@Resolver()
export class AddSiteResolver {

    @Query(() => SiteDto)
    @UseMiddleware(Authentication)
    async addSite(
        @Arg("networkId") networkId: string,
        @Arg("data") data: AddSiteInputDto,
        @Ctx() ctx: Context
    ): Promise<SiteDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.addSite(networkId, data, parseHeaders(ctx));
    }
}
