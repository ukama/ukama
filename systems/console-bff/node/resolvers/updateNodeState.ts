import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { Node, UpdateNodeStateInput } from "./types";

@Resolver()
export class UpdateNodeStateResolver {
  @Mutation(() => Node)
  async updateNodeState(
    @Arg("data") data: UpdateNodeStateInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    return await dataSources.dataSource.updateNodeState({
      id: data.id,
      state: data.state,
    });
  }
}
