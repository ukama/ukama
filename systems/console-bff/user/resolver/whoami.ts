import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { getIdentity } from "../../common/auth/authCalls";
import { WhoamiDto } from "./types";
import { Context } from "../context";

@Resolver()
export class WhoamiResolver {
  @Query(() => WhoamiDto)
  async whoami(@Ctx() ctx: Context): Promise<WhoamiDto> {
    const { dataSources, headers } = ctx;
    return await dataSources.dataSource.whoami(headers.userId);
  }
}
