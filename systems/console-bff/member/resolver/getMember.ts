import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { MemberDto } from "./types";

@Resolver()
export class GetMemberResolver {
  @Query(() => MemberDto)
  async getMember(
    @Arg("id") id: string,
    @Ctx() ctx: Context
  ): Promise<MemberDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getMember(id);
  }
}
