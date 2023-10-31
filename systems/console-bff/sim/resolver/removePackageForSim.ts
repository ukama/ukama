import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import {
  RemovePackageFormSimInputDto,
  RemovePackageFromSimResDto,
} from "./types";

@Resolver()
export class RemovePackageForSimResolver {
  @Mutation(() => RemovePackageFromSimResDto)
  @UseMiddleware(Authentication)
  async removePackageForSim(
    @Arg("data") data: RemovePackageFormSimInputDto,
    @Ctx() ctx: Context
  ): Promise<RemovePackageFromSimResDto> {
    const { dataSources } = ctx;
    return await dataSources.dataSource.removePackageFromSim(data);
  }
}
