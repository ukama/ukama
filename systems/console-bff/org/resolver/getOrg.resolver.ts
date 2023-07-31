import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { Authentication } from "./../../common/auth/index";
import { OrgDto } from "./types";

@Resolver()
export class GetOrgResolver {
  @Query(() => OrgDto)
  @UseMiddleware(Authentication)
  async getOrg(
    @Arg("orgName") orgName: string,
    @Ctx() ctx: Context
  ): Promise<OrgDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getOrg(orgName, parseHeaders());
  }
}
