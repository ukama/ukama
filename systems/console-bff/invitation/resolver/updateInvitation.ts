import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { UpdateInvitationResDto ,UpateInvitationInputDto} from "./types";

@Resolver()
export class UpdateInvitationResolver {
  @Mutation(() => UpdateInvitationResDto)
  @UseMiddleware(Authentication)
  async updateInvitation(
    @Arg("id") id: string,
    @Arg("data") data: UpateInvitationInputDto,
    @Ctx() ctx: Context
  ): Promise<UpdateInvitationResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.updateInvitation(id, data);
  }
}
