import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Context } from "../context";
import { AddNetworkInputDto, NetworkDto } from "../types";
import { Authentication } from "../../common/auth";

@Resolver()
export class AddNetworkResolver {
  @Mutation(() => NetworkDto)
  @UseMiddleware(Authentication)
  async addNetwork(
    @Arg("data") data: AddNetworkInputDto,
    @Ctx() ctx: Context
  ): Promise<NetworkDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.addNetwork(data);
  }
}
