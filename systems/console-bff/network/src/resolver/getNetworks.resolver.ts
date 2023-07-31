import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworksResDto } from "../types";

@Resolver()
export class GetNetworksResolver {

    @Query(() => NetworksResDto)
    @UseMiddleware(Authentication)
    async getNetworks(@Ctx() ctx: Context): Promise<NetworksResDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getNetworks(parseHeaders(ctx));
    }
}
