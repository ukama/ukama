import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { PackageDto, UpdatePackageInputDto } from "./types";

@Resolver()
export class UpdatePackageResolver {
  @Mutation(() => PackageDto)
  @UseMiddleware(Authentication)
  async updatePackage(
    @Arg("packageId") packageId: string,
    @Arg("data") data: UpdatePackageInputDto,
    @Ctx() ctx: Context
  ): Promise<PackageDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.updatePackage(packageId, data);
  }
}
