import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworkDto } from "../types";

@Resolver()
export class GetNetworkResolver {

    @Query(() => NetworkDto)
    @UseMiddleware(Authentication)
    async getNetwork(
        @Arg("networkId") networkId: string,
        @Ctx() ctx: Context
    ): Promise<NetworkDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getNetwork(networkId, parseHeaders(ctx));
    }
}
