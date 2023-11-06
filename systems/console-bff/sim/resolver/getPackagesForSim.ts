import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { GetPackagesForSimInputDto, GetSimPackagesDtoAPI } from "./types";

@Resolver()
export class GetPackagesForSimResolver {
  @Query(() => GetSimPackagesDtoAPI)
  async getPackagesForSim(
    @Arg("data") data: GetPackagesForSimInputDto,
    @Ctx() ctx: Context
  ): Promise<GetSimPackagesDtoAPI> {
    const { dataSources } = ctx;

    return await dataSources.dataSource.getPackagesForSim(data);
  }
}
