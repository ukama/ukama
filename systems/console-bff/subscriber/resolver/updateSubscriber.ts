import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { BoolResponse } from "../../common/types";
import { Context } from "../context";
import { UpdateSubscriberInputDto } from "./types";

@Resolver()
export class UpdateSubscriberResolver {
  @Mutation(() => BoolResponse)
  @UseMiddleware(Authentication)
  async updateSubscriber(
    @Arg("subscriberId") subscriberId: string,
    @Arg("data") data: UpdateSubscriberInputDto,
    @Ctx() ctx: Context
  ): Promise<BoolResponse> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.updateSubscriber(subscriberId, data);
  }
}
