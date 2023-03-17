import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworkService } from "../service";
import { NetworksResDto } from "../types";

@Service()
@Resolver()
export class GetNetworksResolver {
    constructor(private readonly networkService: NetworkService) {}

    @Query(() => NetworksResDto)
    @UseMiddleware(Authentication)
    async getNetworks(@Ctx() ctx: Context): Promise<NetworksResDto> {
        return this.networkService.getNetworks(parseCookie(ctx));
    }
}
