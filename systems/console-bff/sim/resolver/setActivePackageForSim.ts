import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import {
  SetActivePackageForSimInputDto,
  SetActivePackageForSimResDto,
} from "./types";

@Resolver()
export class SetActivePackageForSimResolver {
  @Mutation(() => SetActivePackageForSimResDto)
  async setActivePackageForSim(
    @Arg("data") data: SetActivePackageForSimInputDto,
    @Ctx() ctx: Context
  ): Promise<SetActivePackageForSimResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.setActivePackageForSim(data);
  }
}
