import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimPoolStatsDto } from "../types";

@Resolver()
export class GetSimPoolStatsResolver {

    @Query(() => SimPoolStatsDto)
    @UseMiddleware(Authentication)
    async getSimPoolStats(
        @Arg("type") type: string,
        @Ctx() ctx: Context
    ): Promise<SimPoolStatsDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.getSimPoolStats(type, parseHeaders(ctx));
    }
}
