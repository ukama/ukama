import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { SubscriberDto, SubscriberInputDto } from "./types";

@Resolver()
export class AddSubscriberResolver {
  @Mutation(() => SubscriberDto)
  async addSubscriber(
    @Arg("data") data: SubscriberInputDto,
    @Ctx() ctx: Context
  ): Promise<SubscriberDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.addSubscriber(data);
  }
}
