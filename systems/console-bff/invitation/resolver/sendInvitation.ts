import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { SendInvitationInputDto, SendInvitationResDto } from "./types";

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
