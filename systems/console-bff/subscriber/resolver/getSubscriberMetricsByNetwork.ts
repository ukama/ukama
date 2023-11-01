import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { SubscriberMetricsByNetworkDto } from "./types";

@Resolver()
export class GetSubscriberMetricsByNetworkResolver {
  @Query(() => SubscriberMetricsByNetworkDto)
  async getSubscriberMetricsByNetwork(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: Context
  ): Promise<SubscriberMetricsByNetworkDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSubMetricsByNetwork();
  }
}
