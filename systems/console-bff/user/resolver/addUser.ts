import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { UserInputDto, UserResDto } from "./types";

@Resolver()
export class AddUserResolver {
  @Mutation(() => UserResDto)
  @UseMiddleware(Authentication)
  async addUser(
    @Arg("data") data: UserInputDto,
    @Ctx() ctx: Context
  ): Promise<UserResDto> {
    const { dataSources } = ctx;
    const user = await dataSources.dataSource.addUser(data);
    return user;
  }
}
