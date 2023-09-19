import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { SimPoolStatsDto } from "./types";

@Resolver()
export class GetSimPoolStatsResolver {
  @Query(() => SimPoolStatsDto)
  @UseMiddleware(Authentication)
  async getSimPoolStats(
    @Arg("type") type: string,
    @Ctx() ctx: Context
  ): Promise<SimPoolStatsDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSimPoolStats(type);
  }
}
