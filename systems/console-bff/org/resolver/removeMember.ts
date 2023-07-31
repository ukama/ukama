import { Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { BoolResponse } from "../../common/types";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";

@Resolver()
export class RemoveMemberResolver {
  @Mutation(() => BoolResponse)
  @UseMiddleware(Authentication)
  async removeMember(@Ctx() ctx: Context): Promise<BoolResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.removeMember(parseHeaders());
  }
}
