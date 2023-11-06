import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AllocateSimAPIDto, AllocateSimInputDto } from "./types";

@Resolver()
export class AllocateSimResolver {
  @Mutation(() => AllocateSimAPIDto)
  async allocateSim(
    @Arg("data") data: AllocateSimInputDto,
    @Ctx() ctx: Context
  ): Promise<AllocateSimAPIDto> {
    console.log("Hello Allocate Sim", data);
    const { dataSources } = ctx;
    return await dataSources.dataSource.allocateSim(data);
  }
}
