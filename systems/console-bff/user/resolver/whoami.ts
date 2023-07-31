import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { WhoamiDto } from "./types";

@Resolver()
export class WhoamiResolver {
  @Query(() => WhoamiDto)
  @UseMiddleware(Authentication)
  async whoami(@Ctx() ctx: Context): Promise<WhoamiDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.whoami();
  }
}
