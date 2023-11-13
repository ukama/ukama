/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW, VERSION } from "../../common/configs";
import { CBooleanResponse } from "../../common/types";
import {
  AddNodeInput,
  AddNodeToSiteInput,
  DeleteNode,
  Node,
  NodeInput,
  Nodes,
  UpdateNodeInput,
  UpdateNodeStateInput,
} from "../resolvers/types";
import { AttachNodeInput } from "./../resolvers/types";
import { parseNodeRes, parseNodesRes } from "./mapper";

const NODES = "nodes";
class NodeAPI extends RESTDataSource {
  baseURL = REGISTRY_API_GW;

  async getNode(args: NodeInput): Promise<Node> {
    return this.get(`/${VERSION}/${NODES}/${args.id}`).then(res =>
      parseNodeRes(res.node)
    );
  }
  async getNodes(args: boolean): Promise<Nodes> {
    return this.get(`/${VERSION}/${NODES}?free=${args}`).then(res =>
      parseNodesRes(res)
    );
  }
  async getNodesByNetwork(networkId: string): Promise<Nodes> {
    return this.get(`/${VERSION}/${NODES}/networks/${networkId}`).then(res =>
      parseNodesRes(res)
    );
  }
  async deleteNodeFromOrg(args: NodeInput): Promise<DeleteNode> {
    return this.delete(`/${VERSION}/${NODES}/${args.id}/sites`).then(() =>
      this.delete(`${args.id}`).then(() => {
        return { id: args.id };
      })
    );
  }
  async attachNode(args: AttachNodeInput): Promise<CBooleanResponse> {
    return this.post(`/${VERSION}/${NODES}/${args.parentNode}/attach`, {
      body: {
        anodel: args.anodel,
        anoder: args.anoder,
      },
    }).then(res => (res ? { success: true } : { success: false }));
  }
  async detachhNode(args: NodeInput): Promise<CBooleanResponse> {
    return this.delete(`/${VERSION}/${NODES}/${args.id}/detach`).then(res =>
      res ? { success: true } : { success: false }
    );
  }
  async addNode(args: AddNodeInput): Promise<Node> {
    return this.post(`/${VERSION}/${NODES}/`, {
      body: {
        name: args.name,
        node_id: args.id,
        org_id: args.orgId,
      },
    }).then(res => parseNodeRes(res.node));
  }
  async addNodeToSite(args: AddNodeToSiteInput): Promise<CBooleanResponse> {
    return this.post(`/${VERSION}/${NODES}/${args.nodeId}/sites`, {
      body: {
        net_id: args.networkId,
        site_id: args.siteId,
      },
    }).then(res => (res ? { success: true } : { success: false }));
  }
  async releaseNodeFromSite(args: NodeInput): Promise<CBooleanResponse> {
    return await this.delete(`/${VERSION}/${NODES}/${args.id}/sites`).then(
      res => (res ? { success: true } : { success: false })
    );
  }
  async updateNodeState(args: UpdateNodeStateInput): Promise<Node> {
    return this.patch(
      `/${VERSION}/${NODES}/${args.id}/state/${args.state}`
    ).then(res => parseNodeRes(res));
  }
  async updateNode(args: UpdateNodeInput): Promise<Node> {
    return this.put(`/${VERSION}/${NODES}/${args.id}`, {
      body: {
        name: args.name,
      },
    }).then(res => parseNodeRes(res.node));
  }
}

export default NodeAPI;
