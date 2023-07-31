import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { NetworkStatusDto } from "../types";

@Resolver()
export class GetNetworkStatusResolver {
  @Query(() => NetworkStatusDto)
  async getNetworkStatus(@Ctx() ctx: Context): Promise<NetworkStatusDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getNetworkStatus("ORG_ID");
  }
}
