import { Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { MembersResDto } from "./types";

@Resolver()
export class GetMembersResolver {
  @Query(() => MembersResDto)
  async getMembers(@Ctx() ctx: Context): Promise<MembersResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getMembers();
  }
}
