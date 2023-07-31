import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { OrgUserSimDto, UpdateUserServiceInput } from "./types";

@Resolver()
export class UpdateUserStatusResolver {
  @Mutation(() => OrgUserSimDto)
  @UseMiddleware(Authentication)
  async updateUserStatus(
    @Arg("data") data: UpdateUserServiceInput,
    @Ctx() ctx: Context
  ): Promise<OrgUserSimDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.updateUserStatus(data, parseHeaders());
  }
}
