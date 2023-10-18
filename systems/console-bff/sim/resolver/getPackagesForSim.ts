import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { GetPackagesForSimInputDto, GetPackagesForSimResDto } from "./types";

@Resolver()
export class GetPackagesForSimResolver {
  @Query(() => GetPackagesForSimResDto)
  @UseMiddleware(Authentication)
  async getSim(
    @Arg("data") data: GetPackagesForSimInputDto,
    @Ctx() ctx: Context
  ): Promise<GetPackagesForSimResDto> {
    const { dataSources } = ctx;

    return await dataSources.dataSource.getPackagesForSim(data);
  }
}
