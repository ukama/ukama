import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { UserResDto } from "./types";

@Service()
@Resolver()
export class DeactivateUserResolver {
  @Mutation(() => UserResDto)
  @UseMiddleware(Authentication)
  async deactivateUser(
    @Arg("uuid")
    uuid: string,
    @Ctx() ctx: Context
  ): Promise<UserResDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.deactivateUser(uuid);
  }
}
