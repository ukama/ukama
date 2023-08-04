import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { MemberInputDto } from "./types";

@Resolver()
export class RemoveMemberResolver {
  @Mutation(() => CBooleanResponse)
  async removeMember(
    @Arg("data") data: MemberInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.removeMember(data);
  }
}
