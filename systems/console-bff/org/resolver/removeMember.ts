import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { MemberInputDto } from "./types";

@Resolver()
export class RemoveMemberResolver {
  @Mutation(() => CBooleanResponse)
  @UseMiddleware(Authentication)
  async removeMember(
    @Arg("data") data: MemberInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.removeMember(data);
  }
}
