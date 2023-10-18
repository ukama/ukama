import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { DeleteInvitationResDto } from "./types";

@Resolver()
export class DeleteInvitationResolver {
  @Mutation(() => DeleteInvitationResDto)
  @UseMiddleware(Authentication)
  async deleteInvitation(
    @Arg("id") id: string,
    @Ctx() ctx: Context
  ): Promise<DeleteInvitationResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.deleteInvitation(id);
  }
}
