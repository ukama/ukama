import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { AddOrgInputDto, OrgDto } from "./types";

@Resolver()
export class AddOrgResolver {
  @Mutation(() => OrgDto)
  @UseMiddleware(Authentication)
  async addOrg(
    @Arg("data") data: AddOrgInputDto,
    @Ctx() ctx: Context
  ): Promise<OrgDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.addOrg(data, parseHeaders());
  }
}
