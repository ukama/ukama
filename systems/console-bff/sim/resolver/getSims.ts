import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { SIM_TYPES } from "../../common/enums";
import { Context } from "../context";
import { SimsResDto } from "./types";

@Resolver()
export class GetSimsResolver {
  @Query(() => SimsResDto)
  @UseMiddleware(Authentication)
  async getSims(
    @Arg("type") type: SIM_TYPES,
    @Ctx() ctx: Context
  ): Promise<SimsResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getSims(type);
  }
}
