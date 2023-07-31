import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { AddSiteInputDto, SiteDto } from "../types";
import { Context } from "../context";

@Resolver()
export class AddSiteResolver {
  @Query(() => SiteDto)
  async addSite(
    @Arg("networkId") networkId: string,
    @Arg("data") data: AddSiteInputDto,
    @Ctx() ctx: Context
  ): Promise<SiteDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.addSite(networkId, data);
  }
}
