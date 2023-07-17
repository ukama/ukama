import { Arg, Ctx, Mutation, Query, Resolver } from "type-graphql";

import { Context } from "../context";
import { CBooleanResponse } from "./../../common/types/index";
import {
  AddNodeInput,
  AddNodeToNetworkInput,
  AttachNodeInput,
  DeleteNode,
  GetNode,
  GetNodes,
  Node,
  NodeInput,
  NodeState,
  UpdateNodeInput,
} from "./types";

@Resolver(Node)
class NodeResolvers {
  @Query(() => GetNode)
  async getNode(@Arg("data") data: NodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    const node = await dataSources.dataSource.getNode({ id: data.id });
    return node;
  }

  @Query(() => GetNodes)
  async getNodes(@Ctx() context: Context) {
    const { dataSources } = context;
    const nodes = await dataSources.dataSource.getNodes();
    return nodes;
  }

  @Query(() => GetNodes)
  async getFreeNodes(@Ctx() context: Context) {
    const { dataSources } = context;
    const nodes = await dataSources.dataSource.getFreeNodes();
    return nodes;
  }

  @Mutation(() => DeleteNode)
  async deleteNodeFromOrg(
    @Arg("data") data: NodeInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    const node = await dataSources.dataSource.deleteNodeFromOrg({
      id: data.id,
    });
    return { id: node.node };
  }

  @Mutation(() => CBooleanResponse)
  async attachNode(
    @Arg("data") data: AttachNodeInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    const node = await dataSources.dataSource.attachNode({
      anodel: data.anodel,
      anoder: data.anoder,
      parentNode: data.parentNode,
    });
    return node;
  }

  @Mutation(() => CBooleanResponse)
  async detachhNode(@Arg("data") data: NodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    const node = await dataSources.dataSource.detachhNode({
      id: data.id,
    });
    return node;
  }

  @Mutation(() => GetNode)
  async addNode(@Arg("data") data: AddNodeInput, @Ctx() context: Context) {
    const { dataSources } = context;
    const node = await dataSources.dataSource.addNode({
      id: data.id,
      state: data.state,
    });
    return node;
  }

  @Mutation(() => CBooleanResponse)
  async releaseNodeFromNetwork(
    @Arg("data") data: NodeInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    const node = await dataSources.dataSource.releaseNodeFromNetwork({
      id: data.id,
    });
    return node;
  }

  @Mutation(() => CBooleanResponse)
  async addNodeToNetwork(
    @Arg("data") data: AddNodeToNetworkInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    const node = await dataSources.dataSource.addNodeToNetwork({
      networkId: data.networkId,
      nodeId: data.nodeId,
    });
    return node;
  }

  @Mutation(() => NodeState)
  async updateNodeState(
    @Arg("data") data: UpdateNodeState,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    const node = await dataSources.dataSource.updateNodeState({
      id: data.id,
      state: data.state,
    });
    return node;
  }

  @Mutation(() => GetNode)
  async updateNode(
    @Arg("data") data: UpdateNodeInput,
    @Ctx() context: Context
  ) {
    const { dataSources } = context;
    const node = await dataSources.dataSource.updateNode({
      id: data.id,
      name: data.name,
    });
    return node;
  }
}

export default NodeResolvers;
