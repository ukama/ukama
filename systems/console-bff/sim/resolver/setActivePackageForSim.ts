import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import {
  SetActivePackageForSimInputDto,
  SetActivePackageForSimResDto,
} from "./types";

@Resolver()
export class SetActivePackageForSimResolver {
  @Mutation(() => SetActivePackageForSimResDto)
  @UseMiddleware(Authentication)
  async setActivePackageForSim(
    @Arg("data") data: SetActivePackageForSimInputDto,
    @Ctx() ctx: Context
  ): Promise<SetActivePackageForSimResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.setActivePackageForSim(data);
  }
}
