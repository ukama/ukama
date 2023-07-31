import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { SimDataUsage } from "./types";

@Resolver()
export class GetDataUsageResolver {
  @Query(() => SimDataUsage)
  @UseMiddleware(Authentication)
  async getDataUsage(
    @Arg("simId") simId: string,
    @Ctx() ctx: Context
  ): Promise<SimDataUsage> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getDataUsage(simId);
  }
}
