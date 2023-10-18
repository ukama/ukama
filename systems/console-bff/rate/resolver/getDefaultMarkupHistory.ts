import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { DefaultMarkupHistoryResDto } from "./types";

@Resolver()
export class GetDefaultMarkupHistoryResolver {
  @Query(() => DefaultMarkupHistoryResDto)
  @UseMiddleware(Authentication)
  async getDefaultMarkupHistory(
    @Ctx() ctx: Context
  ): Promise<DefaultMarkupHistoryResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getDefaultMarkupHistory();
  }
}
