import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";

@Resolver()
export class DeleteSubscriberResolver {
  @Mutation(() => CBooleanResponse)
  @UseMiddleware(Authentication)
  async deleteSubscriber(
    @Arg("subscriberId") subscriberId: string,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.deleteSubscriber(subscriberId);
  }
}
