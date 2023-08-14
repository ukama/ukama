import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { GetInvitationByOrgResDto } from "./types";

@Resolver()
export class GetInVitationsByOrgResolver {
  @Query(() => GetInvitationByOrgResDto)
  async getInVitationsByOrg(@Ctx() ctx: Context): Promise<GetInvitationByOrgResDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.getInVitationsByOrg(headers.orgName);
  }
}
