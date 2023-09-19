import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddOrgInputDto, OrgDto } from "./types";

@Resolver()
export class AddOrgResolver {
  @Mutation(() => OrgDto)
  async addOrg(
    @Arg("data") data: AddOrgInputDto,
    @Ctx() ctx: Context
  ): Promise<OrgDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.addOrg(data);
  }
}
