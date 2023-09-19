import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { DefaultMarkupResDto } from "./types";

@Resolver()
export class GetDefaultMarkupResolver {
  @Query(() => DefaultMarkupResDto)
  @UseMiddleware(Authentication)
  async getDefaultMarkup(@Ctx() ctx: Context): Promise<DefaultMarkupResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getDefaultMarkup();
  }
}
