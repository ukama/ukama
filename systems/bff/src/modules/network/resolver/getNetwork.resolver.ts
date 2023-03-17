import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworkService } from "../service";
import { NetworkDto } from "../types";

@Service()
@Resolver()
export class GetNetworkResolver {
    constructor(private readonly networkService: NetworkService) {}

    @Query(() => NetworkDto)
    @UseMiddleware(Authentication)
    async getNetwork(
        @Arg("networkId") networkId: string,
        @Ctx() ctx: Context,
    ): Promise<NetworkDto> {
        return this.networkService.getNetwork(networkId, parseCookie(ctx));
    }
}
