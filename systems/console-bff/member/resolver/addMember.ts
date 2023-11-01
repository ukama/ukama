import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddMemberInputDto, MemberDto } from "./types";

@Resolver()
export class AddMemberResolver {
  @Mutation(() => MemberDto)
  async addMember(
    @Arg("data") data: AddMemberInputDto,
    @Ctx() ctx: Context
  ): Promise<MemberDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.addMember(data);
  }
}
