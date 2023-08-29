import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
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
