import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";

@Resolver()
export class RemoveMemberResolver {
  @Mutation(() => CBooleanResponse)
  async removeMember(
    @Arg("id") id: string,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.removeMember(id);
  }
}
