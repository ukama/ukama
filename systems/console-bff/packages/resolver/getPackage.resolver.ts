import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { PackageDto } from "./types";

@Resolver()
export class GetPackageResolver {
  @Query(() => PackageDto)
  @UseMiddleware(Authentication)
  async getPackage(
    @Arg("packageId") packageId: string,
    @Ctx() ctx: Context
  ): Promise<PackageDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getPackage(packageId);
  }
}
