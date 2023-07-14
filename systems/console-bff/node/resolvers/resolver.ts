import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { NodeInput, NodeRes } from "./types";

@Resolver(NodeRes)
class NodeResolvers {
  @Query(() => NodeRes)
  async getMetrics(@Arg("data") data: NodeInput, @Ctx() context: Context) {
    // const { dataSources } = context;
    // dataSources.dataSource.getNode(data.nodeId);
    return { nodeId: data.nodeId };
  }
}

export default NodeResolvers;
