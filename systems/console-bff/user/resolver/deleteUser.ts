import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { BoolResponse } from "../../common/types";
import { Context } from "../context";

@Resolver()
export class DeleteUserResolver {
  @Mutation(() => BoolResponse)
  @UseMiddleware(Authentication)
  async deleteUser(
    @Arg("userId") userId: string,
    @Ctx() ctx: Context
  ): Promise<BoolResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.deleteUser(userId);
  }
}
