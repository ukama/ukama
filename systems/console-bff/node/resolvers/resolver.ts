import { Arg, Ctx, Mutation, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { CBooleanResponse } from "./../../common/types/index";
import {
  AddNodeInput,
  AddNodeToSiteInput,
  AttachNodeInput,
  DeleteNode,
  GetNode,
  GetNodes,
  GetNodesInput,
  Node,
  NodeInput,
  UpdateNodeInput,
  UpdateNodeStateInput,
} from "./types";

@Resolver(Node)
class NodeResolvers {
  @Query(() => Node)
  async getNode(@Arg("data") data: NodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    return await dataSources.dataSource.getNode({ id: data.id });
  }

  @Query(() => GetNodes)
  async getNodes(@Arg("data") data: GetNodesInput, @Ctx() context: Context) {
    const { dataSources } = context;
    return await dataSources.dataSource.getNodes(data?.isFree || false);
  }

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

  @Mutation(() => CBooleanResponse)
  async detachhNode(@Arg("data") data: NodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    return await dataSources.dataSource.detachhNode({
      id: data.id,
    });
  }

  @Mutation(() => Node)
  async addNode(@Arg("data") data: AddNodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    return await dataSources.dataSource.addNode({
      id: data.id,
      name: data.name,
      orgId: data.orgId,
    });
  }

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

  @Mutation(() => CBooleanResponse)
  async releaseNodeFromSite(
    @Arg("data") data: NodeInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    return await dataSources.dataSource.releaseNodeFromSite({
      id: data.id,
    });
  }

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

export default NodeResolvers;
