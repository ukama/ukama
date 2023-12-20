import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { GetSimBySubscriberInputDto, SubscriberToSimsDto } from "./types";

@Resolver()
export class GetSimsBySubscriberResolver {
  @Query(() => SubscriberToSimsDto)
  @UseMiddleware(Authentication)
  async getSimsBySubscriber(
    @Arg("data") data: GetSimBySubscriberInputDto,
    @Ctx() ctx: Context
  ): Promise<SubscriberToSimsDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSimsBySubscriberId(data);
  }
}
