import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddNetworkInputDto, NetworkDto } from "./types";

@Resolver()
export class AddNetworkResolver {
  @Mutation(() => NetworkDto)
  async addNetwork(
    @Arg("data") data: AddNetworkInputDto,
    @Ctx() ctx: Context
  ): Promise<NetworkDto> {
    const { dataSources } = ctx;
    return dataSources.dataSource.addNetwork(data);
  }
}
