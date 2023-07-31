import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { IdResponse } from "../../common/types";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";

@Resolver()
export class DeletePackageResolver {
  @Mutation(() => IdResponse)
  @UseMiddleware(Authentication)
  async deletePackage(
    @Arg("packageId") packageId: string,
    @Ctx() ctx: Context
  ): Promise<IdResponse> {
    const { dataSources } = ctx;
    return dataSources.dataSource.deletePackage(packageId, parseHeaders());
  }
}
