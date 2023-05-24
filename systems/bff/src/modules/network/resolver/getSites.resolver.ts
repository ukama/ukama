import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworkService } from "../service";
import { SitesResDto } from "../types";

@Service()
@Resolver()
export class GetSitesResolver {
    constructor(private readonly networkService: NetworkService) {}

    @Query(() => SitesResDto)
    @UseMiddleware(Authentication)
    async getSites(
        @Arg("networkId") networkId: string,
        @Ctx() ctx: Context
    ): Promise<SitesResDto> {
        return this.networkService.getSites(networkId, parseHeaders(ctx));
    }
}
