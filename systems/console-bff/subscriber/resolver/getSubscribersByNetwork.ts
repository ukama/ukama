import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { SubscribersResDto } from "./types";

@Resolver()
export class GetSubscribersByNetworkResolver {
  @Query(() => SubscribersResDto)
  async getSubscribersByNetwork(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: Context
  ): Promise<SubscribersResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSubscribersByNetwork(networkId);
  }
}
