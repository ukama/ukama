import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { NodeInput } from "./types";

@Resolver()
export class ReleaseNodeFromSiteResolver {
  @Mutation(() => CBooleanResponse)
  async releaseNodeFromSite(
    @Arg("data") data: NodeInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    return await dataSources.dataSource.releaseNodeFromSite({
      id: data.id,
    });
  }
}
