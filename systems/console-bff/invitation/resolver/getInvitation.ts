import { Ctx, Query, Resolver } from "type-graphql";
import { Arg, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { InvitationDto } from "./types";

@Resolver()
export class GetInvitationResolver {
  @Query(() => InvitationDto)
  @UseMiddleware(Authentication)
  async getInvitation(
    @Arg("id") id: string,
    @Ctx() ctx: Context
  ): Promise<InvitationDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.getInvitation(id);
  }
}
