import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { AllocateSimInputDto, SimDto } from "./types";

@Resolver()
export class AllocateSimResolver {
  @Mutation(() => SimDto)
  @UseMiddleware(Authentication)
  async allocateSim(
    @Arg("data") data: AllocateSimInputDto,
    @Ctx() ctx: Context
  ): Promise<SimDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.allocateSim(data);
  }
}
