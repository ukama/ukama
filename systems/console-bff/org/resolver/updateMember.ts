import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { BoolResponse } from "../../common/types";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { UpdateMemberInputDto } from "./types";

@Resolver()
export class UpdateMemberResolver {
  @Mutation(() => BoolResponse)
  @UseMiddleware(Authentication)
  async updateMember(
    @Arg("memberId") memberId: string,
    @Arg("data") data: UpdateMemberInputDto,
    @Ctx() ctx: Context
  ): Promise<BoolResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.updateMember(memberId, data, parseHeaders());
  }
}
