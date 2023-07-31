import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { UserFistVisitInputDto, UserFistVisitResDto } from "./types";

@Resolver()
export class updateFirstVisitResolver {
  @Mutation(() => UserFistVisitResDto)
  @UseMiddleware(Authentication)
  async updateFirstVisit(
    @Arg("data") data: UserFistVisitInputDto,
    @Ctx() ctx: Context
  ): Promise<UserFistVisitResDto> {
    const { dataSources } = ctx;
    const user = await dataSources.dataSource.updateFirstVisit(
      data,
      parseHeaders()
    );
    return user;
  }
}
