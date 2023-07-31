import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { UserResDto } from "./types";

@Resolver()
export class GetUserResolver {
  @Query(() => UserResDto)
  @UseMiddleware(Authentication)
  async getUser(
    @Arg("userId") userId: string,
    @Ctx() ctx: Context
  ): Promise<UserResDto | null> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getUser(userId);
  }
}
