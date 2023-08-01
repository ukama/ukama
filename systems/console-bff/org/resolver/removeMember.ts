import { Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { CBooleanResponse } from "../../common/types";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";

@Resolver()
export class RemoveMemberResolver {
  @Mutation(() => CBooleanResponse)
  @UseMiddleware(Authentication)
  async removeMember(@Ctx() ctx: Context): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.removeMember(parseHeaders());
  }
}
