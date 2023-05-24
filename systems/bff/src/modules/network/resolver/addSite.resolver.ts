import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworkService } from "../service";
import { AddSiteInputDto, SiteDto } from "./../types";

@Service()
@Resolver()
export class addSiteResolver {
    constructor(private readonly networkService: NetworkService) {}

    @Query(() => SiteDto)
    @UseMiddleware(Authentication)
    async addSite(
        @Arg("networkId") networkId: string,
        @Arg("data") data: AddSiteInputDto,
        @Ctx() ctx: Context
    ): Promise<SiteDto> {
        return this.networkService.addSite(networkId, data, parseHeaders(ctx));
    }
}
