import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { AddPackageSimResDto, AddPackageToSimInputDto } from "./types";

@Resolver()
export class AddPackageToSimResolver {
  @Mutation(() => AddPackageSimResDto)
  @UseMiddleware(Authentication)
  async addPackageToSim(
    @Arg("data") data: AddPackageToSimInputDto,
    @Ctx() ctx: Context
  ): Promise<AddPackageSimResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.addPackageToSim(data);
  }
}
