import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { MemberObj } from "./types";

@Resolver()
export class GetOrgMemberResolver {
  @Query(() => MemberObj)
  @UseMiddleware(Authentication)
  async getOrgMember(@Ctx() ctx: Context): Promise<MemberObj> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getOrgMember(parseHeaders());
  }
}
