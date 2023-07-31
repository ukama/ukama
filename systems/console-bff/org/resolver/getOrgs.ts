import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { OrgsResDto } from "./types";

@Resolver()
export class GetOrgsResolver {
  @Query(() => OrgsResDto)
  @UseMiddleware(Authentication)
  async getOrgs(@Ctx() ctx: Context): Promise<OrgsResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getOrgs(parseHeaders());
  }
}
