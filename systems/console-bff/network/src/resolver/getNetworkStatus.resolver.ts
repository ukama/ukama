import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NetworkStatusDto } from "../types";

@Resolver()
export class GetNetworkStatusResolver {

    @Query(() => NetworkStatusDto)
    @UseMiddleware(Authentication)
    async getNetworkStatus(@Ctx() ctx: Context): Promise<NetworkStatusDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getNetworkStatus(parseHeaders(ctx));
    }
}
