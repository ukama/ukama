import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { UpdateSubscriberInputDto } from "./types";

@Resolver()
export class UpdateSubscriberResolver {
  @Mutation(() => CBooleanResponse)
  async updateSubscriber(
    @Arg("subscriberId") subscriberId: string,
    @Arg("data") data: UpdateSubscriberInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.updateSubscriber(subscriberId, data);
  }
}
