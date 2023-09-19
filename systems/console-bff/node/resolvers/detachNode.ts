import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { NodeInput } from "./types";

@Resolver()
export class DetachNodeResolver {
  @Mutation(() => CBooleanResponse)
  async detachhNode(@Arg("data") data: NodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    return await dataSources.dataSource.detachhNode({
      id: data.id,
    });
  }
}
