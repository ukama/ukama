import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { AddMemberInputDto, MemberObj } from "./types";

@Resolver()
export class AddMemberResolver {
  @Mutation(() => MemberObj)
  @UseMiddleware(Authentication)
  async addMember(
    @Arg("data") data: AddMemberInputDto,
    @Ctx() ctx: Context
  ): Promise<MemberObj> {
    const { dataSources } = ctx;
    return dataSources.dataSource.addMember(data, parseHeaders());
  }
}
