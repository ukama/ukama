import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { AddPackageInputDto, PackageDto } from "./types";

@Resolver()
export class AddPackageResolver {
  @Mutation(() => PackageDto)
  @UseMiddleware(Authentication)
  async addPackage(
    @Arg("data") data: AddPackageInputDto,
    @Ctx() ctx: Context
  ): Promise<PackageDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.addPackage(data, headers);
  }
}
