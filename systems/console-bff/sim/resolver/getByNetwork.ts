import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { GetSimByNetworkInputDto, SimDetailsDto } from "./types";

@Resolver()
export class GetSimByNetworkResolver {
  @Query(() => SimDetailsDto)
  @UseMiddleware(Authentication)
  async getSim(
    @Arg("data") data: GetSimByNetworkInputDto,
    @Ctx() ctx: Context
  ): Promise<SimDetailsDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getSimByNetworkId(data);
  }
}
