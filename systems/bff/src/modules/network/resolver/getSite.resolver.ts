import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworkService } from "../service";
import { SiteDto } from "./../types";

@Service()
@Resolver()
export class GetSiteResolver {
    constructor(private readonly networkService: NetworkService) {}

    @Query(() => SiteDto)
    @UseMiddleware(Authentication)
    async getSite(
        @Arg("networkId") networkId: string,
        @Arg("siteId") siteId: string,
        @Ctx() ctx: Context
    ): Promise<SiteDto> {
        return this.networkService.getSite(
            networkId,
            siteId,
            parseHeaders(ctx)
        );
    }
}
