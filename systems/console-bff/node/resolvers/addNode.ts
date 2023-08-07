import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { Context } from "../context";
import { AddNodeInput, Node } from "./types";

@Resolver()
export class AddNodeResolver {
  @Mutation(() => Node)
  async addNode(@Arg("data") data: AddNodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    return await dataSources.dataSource.addNode({
      id: data.id,
      name: data.name,
      orgId: data.orgId,
    });
  }
}
