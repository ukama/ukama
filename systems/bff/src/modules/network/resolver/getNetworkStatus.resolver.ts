import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworkService } from "../service";
import { NetworkStatusDto } from "../types";

@Service()
@Resolver()
export class GetNetworkStatusResolver {
    constructor(private readonly networkService: NetworkService) {}

    @Query(() => NetworkStatusDto)
    @UseMiddleware(Authentication)
    async getNetworkStatus(@Ctx() ctx: Context): Promise<NetworkStatusDto> {
        return this.networkService.getNetworkStatus(parseCookie(ctx));
    }
}
