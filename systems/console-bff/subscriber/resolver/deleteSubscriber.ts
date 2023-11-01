import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";

@Resolver()
export class DeleteSubscriberResolver {
  @Mutation(() => CBooleanResponse)
  async deleteSubscriber(
    @Arg("subscriberId") subscriberId: string,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.deleteSubscriber(subscriberId);
  }
}
