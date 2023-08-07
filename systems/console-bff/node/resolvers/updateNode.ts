import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { Node, UpdateNodeInput } from "./types";

@Resolver()
export class UpdateNodeResolver {
  @Mutation(() => Node)
  async updateNode(
    @Arg("data") data: UpdateNodeInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    return await dataSources.dataSource.updateNode({
      id: data.id,
      name: data.name,
    });
  }
}
