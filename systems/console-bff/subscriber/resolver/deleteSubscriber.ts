import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { BoolResponse } from "../../common/types";
import { Context } from "../context";

@Resolver()
export class DeleteSubscriberResolver {
  @Mutation(() => BoolResponse)
  @UseMiddleware(Authentication)
  async deleteSubscriber(
    @Arg("subscriberId") subscriberId: string,
    @Ctx() ctx: Context
  ): Promise<BoolResponse> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.deleteSubscriber(subscriberId);
  }
}
