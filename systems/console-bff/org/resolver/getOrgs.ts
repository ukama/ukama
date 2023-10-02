import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { OrgsResDto } from "./types";

@Resolver()
export class GetOrgsResolver {
  @Query(() => OrgsResDto)
  async getOrgs(@Ctx() ctx: Context): Promise<OrgsResDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.getOrgs(headers.userId);
  }
}
