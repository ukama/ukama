import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddPackageInputDto, PackageDto } from "./types";

@Resolver()
export class AddPackageResolver {
  @Mutation(() => PackageDto)
  async addPackage(
    @Arg("data") data: AddPackageInputDto,
    @Ctx() ctx: Context
  ): Promise<PackageDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.addPackage(data, headers);
  }
}
