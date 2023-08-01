import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { OrgDto } from "./types";

@Resolver()
export class GetOrgResolver {
  @Query(() => OrgDto)
  async getOrg(
    @Arg("orgName") orgName: string,
    @Ctx() ctx: Context
  ): Promise<OrgDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getOrg(orgName);
  }
}
