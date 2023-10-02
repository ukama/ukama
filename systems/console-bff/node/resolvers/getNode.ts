import { Arg, Ctx, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { Node, NodeInput } from "./types";

@Resolver()
export class GetNodeResolver {
  @Query(() => Node)
  async getNode(@Arg("data") data: NodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    return await dataSources.dataSource.getNode({ id: data.id });
  }
}
