import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetSimBySubscriberInputDto, SubscriberToSimsDto } from "./types";

@Resolver()
export class GetSimsBySubscriberResolver {
  @Query(() => SubscriberToSimsDto)
  async getSimsBySubscriber(
    @Arg("data") data: GetSimBySubscriberInputDto,
    @Ctx() ctx: AppContext
  ): Promise<SubscriberToSimsDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("sim");
    return await dataSources.sim.getSimsBySubscriberId(baseURL, data);
  }
}
