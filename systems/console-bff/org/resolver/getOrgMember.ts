import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { MemberInputDto, MemberObj } from "./types";

@Resolver()
export class GetOrgMemberResolver {
  @Query(() => MemberObj)
  async getOrgMember(
    @Arg("data") data: MemberInputDto,
    @Ctx() ctx: Context
  ): Promise<MemberObj> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getOrgMember(data);
    //TODO: MAKE GET USER BY ID CALL HERE
  }
}
