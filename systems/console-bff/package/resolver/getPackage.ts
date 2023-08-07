import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { PackageDto } from "./types";

@Resolver()
export class GetPackageResolver {
  @Query(() => PackageDto)
  async getPackage(
    @Arg("packageId") packageId: string,
    @Ctx() ctx: Context
  ): Promise<PackageDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getPackage(packageId);
  }
}
