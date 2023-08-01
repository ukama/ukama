import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddSiteInputDto, SiteDto } from "./types";

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
