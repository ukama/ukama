/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import { CBooleanResponse } from "../../common/types";
import {
  AddNodeInput,
  AddNodeToSiteInput,
  DeleteNode,
  GetNodesByStateInput,
  Node,
  NodeInput,
  NodeStateRes,
  Nodes,
  NodesFilterInput,
  UpdateNodeInput,
  UpdateNodeStateInput,
} from "../resolvers/types";
import { AttachNodeInput } from "./../resolvers/types";
import { getNodeState, parseNodeRes, parseNodesRes } from "./mapper";

const NODES = "nodes";
class NodeAPI extends RESTDataSource {
  async getNode(baseURL: string, args: NodeInput): Promise<Node> {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${NODES}/${args.id}`).then(res =>
      parseNodeRes(res.node)
    );
  }
  async getNodes(baseURL: string, args: NodesFilterInput): Promise<Nodes> {
    this.baseURL = baseURL;
    this.logger.info(`Getting nodes with args: ${JSON.stringify(args)}`);
    const query = new URLSearchParams();
    if (args.id) query.set("id", args.id);
    if (args.siteId) query.set("siteId", args.siteId);
    if (args.networkId) query.set("networkId", args.networkId);
    if (args.type) query.set("type", args.type.toLowerCase());
    if (args.state) query.set("state", args.state.toLowerCase());
    if (args.connectivity)
      query.set("connectivity", args.connectivity.toLowerCase());
    return this.get(`/${VERSION}/${NODES}/list?${query.toString()}`).then(res =>
      parseNodesRes(res)
    );
  }

  async getNodeState(baseURL: string, id: string): Promise<NodeStateRes> {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/state/${id}/history`).then(res =>
      getNodeState(res)
    );
  }
  async getNodesByState(
    baseURL: string,
    data: GetNodesByStateInput
  ): Promise<Nodes> {
    this.baseURL = baseURL;
    return this.get(
      `/${VERSION}/${NODES}/state?connectivity=${data.connectivity.toLocaleLowerCase()}&state=${data.state.toLocaleLowerCase()}`
    ).then(res => parseNodesRes(res));
  }
  async getNodesByNetwork(baseURL: string, networkId: string): Promise<Nodes> {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${NODES}/networks/${networkId}`).then(res =>
      parseNodesRes(res)
    );
  }
  async deleteNodeFromOrg(
    baseURL: string,
    args: NodeInput
  ): Promise<DeleteNode> {
    this.baseURL = baseURL;
    return this.delete(`/${VERSION}/${NODES}/${args.id}/sites`).then(() =>
      this.delete(`${args.id}`).then(() => {
        return { id: args.id };
      })
    );
  }
  async attachNode(
    baseURL: string,
    args: AttachNodeInput
  ): Promise<CBooleanResponse> {
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${NODES}/${args.parentNode}/attach`, {
      body: {
        anodel: args.anodel,
        anoder: args.anoder,
      },
    }).then(res => (res ? { success: true } : { success: false }));
  }
  async detachhNode(
    baseURL: string,
    args: NodeInput
  ): Promise<CBooleanResponse> {
    this.baseURL = baseURL;
    return this.delete(`/${VERSION}/${NODES}/${args.id}/detach`).then(res =>
      res ? { success: true } : { success: false }
    );
  }
  async addNode(baseURL: string, args: AddNodeInput): Promise<Node> {
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${NODES}/`, {
      body: {
        name: args.name,
        node_id: args.id,
      },
    }).then(res => parseNodeRes(res.node));
  }
  async addNodeToSite(
    baseURL: string,
    args: AddNodeToSiteInput
  ): Promise<CBooleanResponse> {
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${NODES}/${args.nodeId}/sites`, {
      body: {
        net_id: args.networkId,
        site_id: args.siteId,
      },
    }).then(res => (res ? { success: true } : { success: false }));
  }
  async releaseNodeFromSite(
    baseURL: string,
    args: NodeInput
  ): Promise<CBooleanResponse> {
    this.baseURL = baseURL;
    return await this.delete(`/${VERSION}/${NODES}/${args.id}/sites`).then(
      res => (res ? { success: true } : { success: false })
    );
  }
  async updateNodeState(
    baseURL: string,
    args: UpdateNodeStateInput
  ): Promise<Node> {
    this.baseURL = baseURL;
    return this.patch(
      `/${VERSION}/${NODES}/${args.id}/state/${args.state}`
    ).then(res => parseNodeRes(res));
  }
  async updateNode(baseURL: string, args: UpdateNodeInput): Promise<Node> {
    this.baseURL = baseURL;
    return this.put(`/${VERSION}/${NODES}/${args.id}`, {
      body: {
        name: args.name,
      },
    }).then(res => parseNodeRes(res.node));
  }
}

export default NodeAPI;
