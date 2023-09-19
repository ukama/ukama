import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW } from "../../common/configs";
import { CBooleanResponse } from "../../common/types";
import {
  AddNodeInput,
  AddNodeToSiteInput,
  DeleteNode,
  GetNodes,
  Node,
  NodeInput,
  UpdateNodeInput,
  UpdateNodeStateInput,
} from "../resolvers/types";
import { AttachNodeInput } from "./../resolvers/types";
import { parseNodeRes, parseNodesRes } from "./mapper";

const version = "/v1/nodes/";
class NodeAPI extends RESTDataSource {
  baseURL = REGISTRY_API_GW + version;

  async getNode(args: NodeInput): Promise<Node> {
    return this.get(args.id).then(res => parseNodeRes(res.node));
  }
  async getNodes(args: boolean): Promise<GetNodes> {
    return this.get(`?free=${args}`).then(res => parseNodesRes(res));
  }
  async deleteNodeFromOrg(args: NodeInput): Promise<DeleteNode> {
    return this.delete(`${args.id}/sites`).then(() =>
      this.delete(`${args.id}`).then(() => {
        return { id: args.id };
      })
    );
  }
  async attachNode(args: AttachNodeInput): Promise<CBooleanResponse> {
    return this.post(`${args.parentNode}/attach`, {
      body: {
        anodel: args.anodel,
        anoder: args.anoder,
      },
    }).then(res => (res ? { success: true } : { success: false }));
  }
  async detachhNode(args: NodeInput): Promise<CBooleanResponse> {
    return this.delete(`${args.id}/detach`).then(res =>
      res ? { success: true } : { success: false }
    );
  }
  async addNode(args: AddNodeInput): Promise<Node> {
    return this.post("", {
      body: {
        name: args.name,
        node_id: args.id,
        org_id: args.orgId,
      },
    }).then(res => parseNodeRes(res.node));
  }
  async addNodeToSite(args: AddNodeToSiteInput): Promise<CBooleanResponse> {
    return this.post(`${args.nodeId}/sites`, {
      body: {
        net_id: args.networkId,
        site_id: args.siteId,
      },
    }).then(res => (res ? { success: true } : { success: false }));
  }
  async releaseNodeFromSite(args: NodeInput): Promise<CBooleanResponse> {
    return await this.delete(`${args.id}/sites`).then(res =>
      res ? { success: true } : { success: false }
    );
  }
  async updateNodeState(args: UpdateNodeStateInput): Promise<Node> {
    return this.patch(`${args.id}/state/${args.state}`).then(res =>
      parseNodeRes(res)
    );
  }
  async updateNode(args: UpdateNodeInput): Promise<Node> {
    return this.put(`${args.id}`, {
      body: {
        name: args.name,
      },
    }).then(res => parseNodeRes(res.node));
  }
}

export default NodeAPI;
