import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { PackagesResDto } from "./types";

@Resolver()
export class GetPackagesResolver {
  @Query(() => PackagesResDto)
  @UseMiddleware(Authentication)
  async getPackages(@Ctx() ctx: Context): Promise<PackagesResDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.getPackages(headers);
  }
}
