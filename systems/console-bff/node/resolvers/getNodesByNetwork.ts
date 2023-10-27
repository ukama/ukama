import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { Nodes } from "./types";

@Resolver()
export class GetNodesByNetworkResolver {
  @Query(() => Nodes)
  async getNodesByNetwork(
    @Arg("networkId") networkId: string,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    return await dataSources.dataSource.getNodesByNetwork(networkId);
  }
}
