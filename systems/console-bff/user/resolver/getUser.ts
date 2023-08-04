import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { Context } from "../context";
import { UserResDto } from "./types";

@Resolver()
export class GetUserResolver {
  @Query(() => UserResDto)
  async getUser(
    @Arg("userId") userId: string,
    @Ctx() ctx: Context
  ): Promise<UserResDto | null> {
    const { dataSources, headers } = ctx;
    logger.info("Headers", headers);
    return dataSources.dataSource.getUser(userId);
  }
}
