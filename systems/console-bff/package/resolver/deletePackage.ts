import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { IdResponse } from "../../common/types";
import { Context } from "../context";

@Resolver()
export class DeletePackageResolver {
  @Mutation(() => IdResponse)
  async deletePackage(
    @Arg("packageId") packageId: string,
    @Ctx() ctx: Context
  ): Promise<IdResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.deletePackage(packageId);
  }
}
