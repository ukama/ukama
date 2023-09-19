import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { SiteDto } from "./types";

@Resolver()
export class GetSiteResolver {
  @Query(() => SiteDto)
  async getSite(
    @Arg("networkId") networkId: string,
    @Arg("siteId") siteId: string,
    @Ctx() ctx: Context
  ): Promise<SiteDto> {
    const { dataSources } = ctx;

    return dataSources.dataSource.getSite(networkId, siteId);
  }
}
