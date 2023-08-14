
import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Context } from "../context";
import { SendInvitationInputDto, SendInvitationResDto } from "./types";
import { Authentication } from "../../common/auth";

@Resolver()
export class SendInvitationResolver {
  @Mutation(() => SendInvitationResDto)
  @UseMiddleware(Authentication)
  async sendInvitation(
    @Arg("data") data: SendInvitationInputDto,
    @Ctx() ctx: Context
  ): Promise<SendInvitationResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.sendInvitation(data);
  }
}
