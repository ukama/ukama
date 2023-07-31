import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimDataUsage } from "../types";

@Resolver()
export class GetDataUsageResolver {

    @Query(() => SimDataUsage)
    @UseMiddleware(Authentication)
    async getDataUsage(
        @Arg("simId") simId: string,
        @Ctx() ctx: Context
    ): Promise<SimDataUsage> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.getDataUsage(simId, parseHeaders(ctx));
    }
}
