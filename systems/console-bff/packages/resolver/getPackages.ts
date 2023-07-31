import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { PackagesResDto } from "./types";

@Resolver()
export class GetPackagesResolver {
  @Query(() => PackagesResDto)
  @UseMiddleware(Authentication)
  async getPackages(@Ctx() ctx: Context): Promise<PackagesResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getPackages(parseHeaders());
  }
}
