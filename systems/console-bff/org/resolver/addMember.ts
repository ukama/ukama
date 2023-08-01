import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddMemberInputDto, MemberObj } from "./types";

@Resolver()
export class AddMemberResolver {
  @Mutation(() => MemberObj)
  async addMember(
    @Arg("data") data: AddMemberInputDto,
    @Ctx() ctx: Context
  ): Promise<MemberObj> {
    const { dataSources } = ctx;
    return dataSources.dataSource.addMember(data);
  }
}
