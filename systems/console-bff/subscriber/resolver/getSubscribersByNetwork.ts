import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { SubscribersResDto } from "./types";

@Resolver()
export class GetSubscribersByNetworkResolver {
  @Query(() => SubscribersResDto)
  @UseMiddleware(Authentication)
  async getSubscribersByNetwork(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: Context
  ): Promise<SubscribersResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSubscribersByNetwork(networkId);
  }
}
