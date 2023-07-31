import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { NetworksResDto } from "../types";

@Resolver()
export class GetNetworksResolver {
  @Query(() => NetworksResDto)
  async getNetworks(@Ctx() ctx: Context): Promise<NetworksResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getNetworks("ORG_ID");
  }
}
