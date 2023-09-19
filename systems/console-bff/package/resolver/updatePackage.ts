import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { PackageDto, UpdatePackageInputDto } from "./types";

@Resolver()
export class UpdatePackageResolver {
  @Mutation(() => PackageDto)
  async updatePackage(
    @Arg("packageId") packageId: string,
    @Arg("data") data: UpdatePackageInputDto,
    @Ctx() ctx: Context
  ): Promise<PackageDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.updatePackage(packageId, data);
  }
}
