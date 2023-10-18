import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { PackagesResDto } from "./types";

@Resolver()
export class GetPackagesResolver {
  @Query(() => PackagesResDto)
  async getPackages(@Ctx() ctx: Context): Promise<PackagesResDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.getPackages(headers);
  }
}
