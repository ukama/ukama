import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { GetSimInputDto, SimDto } from "./types";

@Resolver()
export class GetSimResolver {
  @Query(() => SimDto)
  @UseMiddleware(Authentication)
  async getSim(
    @Arg("data") data: GetSimInputDto,
    @Ctx() ctx: Context
  ): Promise<SimDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSim(data);
  }
}
