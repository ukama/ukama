import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { UpdateUserInputDto, UserResDto } from "./types";

@Service()
@Resolver()
export class UpdateUserResolver {
  @Mutation(() => UserResDto)
  @UseMiddleware(Authentication)
  async updateUser(
    @Arg("userId") userId: string,
    @Arg("data") data: UpdateUserInputDto,
    @Ctx() ctx: Context
  ): Promise<UserResDto | null> {
    const { dataSources } = ctx;
    return dataSources.dataSource.updateUser(userId, data);
  }
}
