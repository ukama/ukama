import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AllocateSimInputDto, SimDto } from "./types";

@Resolver()
export class AllocateSimResolver {
  @Mutation(() => SimDto)
  async allocateSim(
    @Arg("data") data: AllocateSimInputDto,
    @Ctx() ctx: Context
  ): Promise<SimDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.allocateSim(data);
  }
}
