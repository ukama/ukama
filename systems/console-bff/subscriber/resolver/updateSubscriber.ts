import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { UpdateSubscriberInputDto } from "./types";

@Resolver()
export class UpdateSubscriberResolver {
  @Mutation(() => CBooleanResponse)
  @UseMiddleware(Authentication)
  async updateSubscriber(
    @Arg("subscriberId") subscriberId: string,
    @Arg("data") data: UpdateSubscriberInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.updateSubscriber(subscriberId, data);
  }
}
