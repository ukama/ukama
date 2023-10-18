import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { SitesResDto } from "./types";

@Resolver()
export class GetSitesResolver {
  @Query(() => SitesResDto)
  async getSites(
    @Arg("networkId") networkId: string,
    @Ctx() ctx: Context
  ): Promise<SitesResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getSites(networkId);
  }
}
