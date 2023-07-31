import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { NetworkDto } from "../types";

@Resolver()
export class GetNetworkResolver {
  @Query(() => NetworkDto)
  async getNetwork(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: Context
  ): Promise<NetworkDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getNetwork(networkId);
  }
}
