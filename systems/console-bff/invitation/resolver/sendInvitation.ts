
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { SendInvitationInputDto, SendInvitationResDto } from "./types";

@Resolver()
export class SendInvitationResolver {
  @Mutation(() => SendInvitationResDto)
  async sendInvitation(
    @Arg("data") data: SendInvitationInputDto,
    @Ctx() ctx: Context
  ): Promise<SendInvitationResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.sendInvitation(data);
  }
}
