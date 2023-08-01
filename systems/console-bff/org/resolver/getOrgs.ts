import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { OrgsResDto } from "./types";

@Resolver()
export class GetOrgsResolver {
  @Query(() => OrgsResDto)
  async getOrgs(@Ctx() ctx: Context): Promise<OrgsResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getOrgs(
      "08a594d7-a292-43cf-9652-54785b03f48f"
    );
  }
}
