import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { CBooleanResponse } from "../../common/types";
import { Context } from "../context";
import { AddNodeToSiteInput } from "./types";

@Resolver()
export class AddNodeToSiteResolver {
  @Mutation(() => CBooleanResponse)
  async addNodeToSite(
    @Arg("data") data: AddNodeToSiteInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    return dataSources.dataSource.addNodeToSite({
      nodeId: data.nodeId,
      networkId: data.networkId,
      siteId: data.siteId,
    });
  }
}
