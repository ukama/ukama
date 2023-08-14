import { Ctx, Query, Resolver } from "type-graphql";
import { Context } from "../context";
import { GetInvitationInputDto } from "./types";

@Resolver()
export class GetInvitationResolver {
  @Query(() => GetInvitationInputDto)
  async getInvitation(@Ctx() ctx: Context): Promise<GetInvitationInputDto> {
    const { dataSources, headers } = ctx;
    return dataSources.dataSource.getInvitation(headers.orgName);
  }
}


