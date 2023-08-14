import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Context } from "../context";
import { GetInvitationByOrgResDto } from "./types";
import { Authentication } from "../../common/auth";

@Resolver()
export class GetInVitationsByOrgResolver {
  @Query(() => GetInvitationByOrgResDto)
  @UseMiddleware(Authentication)
  async getInVitationsByOrg(@Ctx() ctx: Context): Promise<GetInvitationByOrgResDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.getInVitationsByOrg(headers.orgName);
  }
}
