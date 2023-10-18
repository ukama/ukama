import { Ctx, Query, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { Context } from "../context";
import { OrgDto } from "./types";

@Resolver()
export class GetOrgResolver {
  @Query(() => OrgDto)
  async getOrg(@Ctx() ctx: Context): Promise<OrgDto> {
    const { dataSources, headers } = ctx;
    logger.info(`getOrg: ${JSON.stringify(headers)}`);
    return dataSources.dataSource.getOrg(headers.orgName);
  }
}
