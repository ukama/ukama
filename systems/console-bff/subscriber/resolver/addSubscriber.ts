import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { SubscriberDto, SubscriberInputDto } from "./types";

@Resolver()
export class AddSubscriberResolver {
  @Mutation(() => SubscriberDto)
  @UseMiddleware(Authentication)
  async addSubscriber(
    @Arg("data") data: SubscriberInputDto,
    @Ctx() ctx: Context
  ): Promise<SubscriberDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.addSubscriber(data);
  }
}
