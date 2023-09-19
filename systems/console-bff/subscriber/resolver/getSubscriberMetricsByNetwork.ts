import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { SubscriberMetricsByNetworkDto } from "./types";

@Resolver()
export class GetSubscriberMetricsByNetworkResolver {
  @Query(() => SubscriberMetricsByNetworkDto)
  @UseMiddleware(Authentication)
  async getSubscriberMetricsByNetwork(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: Context
  ): Promise<SubscriberMetricsByNetworkDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSubMetricsByNetwork();
  }
}
