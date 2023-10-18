import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { WhoamiDto } from "./types";

@Resolver()
export class WhoamiResolver {
  @Query(() => WhoamiDto)
  async whoami(@Ctx() ctx: Context): Promise<WhoamiDto> {
    const { dataSources, headers } = ctx;
    return await dataSources.dataSource.whoami(headers.userId);
  }
}
