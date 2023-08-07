import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { DeleteNode, NodeInput } from "./types";

@Resolver()
export class DeleteNodeFromOrgResolver {
  @Mutation(() => DeleteNode)
  async deleteNodeFromOrg(
    @Arg("data") data: NodeInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    return await dataSources.dataSource.deleteNodeFromOrg({
      id: data.id,
    });
  }
}
