import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { UserResDto } from "./types";

@Resolver()
export class GetUserResolver {
  @Query(() => UserResDto)
  async getUser(
    @Arg("userId") userId: string,
    @Ctx() ctx: Context
  ): Promise<UserResDto | null> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getUser(userId);
  }
}
