import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { SubscriberDto } from "./types";

@Resolver()
export class GetSubscriberResolver {
  @Query(() => SubscriberDto)
  async getSubscriber(
    @Arg("subscriberId") subscriberId: string,
    @Ctx() ctx: Context
  ): Promise<SubscriberDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSubscriber(subscriberId);
  }
}
