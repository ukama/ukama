import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { GetSimBySubscriberInputDto, SubscriberToSimsDto } from "./types";

@Resolver()
export class GetSimsBySubscriberResolver {
  @Query(() => SubscriberToSimsDto)
  async getSimsBySubscriber(
    @Arg("data") data: GetSimBySubscriberInputDto,
    @Ctx() ctx: Context
  ): Promise<SubscriberToSimsDto> {
    const { dataSources, baseURL } = ctx;
    return await dataSources.dataSource.getSimsBySubscriberId(baseURL, data);
  }
}
