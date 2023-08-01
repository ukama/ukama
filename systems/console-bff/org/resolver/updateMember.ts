import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { CBooleanResponse } from "../../common/types";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { UpdateMemberInputDto } from "./types";

@Resolver()
export class UpdateMemberResolver {
  @Mutation(() => CBooleanResponse)
  async updateMember(
    @Arg("memberId") memberId: string,
    @Arg("data") data: UpdateMemberInputDto,
    @Ctx() ctx: Context
  ): Promise<CBooleanResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.updateMember(memberId, data);
  }
}
