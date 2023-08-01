import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { MemberInputDto, MemberObj } from "./types";

@Resolver()
export class GetOrgMemberResolver {
  @Query(() => MemberObj)
  @UseMiddleware(Authentication)
  async getOrgMember(
    @Arg("data") data: MemberInputDto,
    @Ctx() ctx: Context
  ): Promise<MemberObj> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getOrgMember(data);
    //TODO: MAKE GET USER BY ID CALL HERE
  }
}
