import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { GetNodesInput, Nodes } from "./types";

@Resolver()
export class GetNodesResolver {
  @Query(() => Nodes)
  async getNodes(@Arg("data") data: GetNodesInput, @Ctx() context: Context) {
    const { dataSources } = context;
    return await dataSources.dataSource.getNodes(data?.isFree || false);
  }
}
