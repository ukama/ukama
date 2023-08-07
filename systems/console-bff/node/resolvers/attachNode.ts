import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { AttachNodeInput } from "./types";

@Resolver()
export class AttachNodeResolver {
  @Mutation(() => CBooleanResponse)
  async attachNode(
    @Arg("data") data: AttachNodeInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    return await dataSources.dataSource.attachNode({
      anodel: data.anodel,
      anoder: data.anoder,
      parentNode: data.parentNode,
    });
  }
}
