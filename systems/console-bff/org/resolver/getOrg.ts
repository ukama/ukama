import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { OrgDto } from "./types";

@Resolver()
export class GetOrgResolver {
  @Query(() => OrgDto)
  async getOrg(@Ctx() ctx: Context): Promise<OrgDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.getOrg(headers.orgName);
  }
}
