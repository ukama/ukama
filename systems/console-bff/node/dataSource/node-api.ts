import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

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
    return this.get(`/${VERSION}/${NODES}/${args.id}`)
      .then(res => parseNodeRes(res.node))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async getNodes(args: boolean): Promise<Nodes> {
    return this.get(`/${VERSION}/${NODES}?free=${args}`)
      .then(res => parseNodesRes(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async getNodesByNetwork(networkId: string): Promise<Nodes> {
    return this.get(`/${VERSION}/${NODES}/networks/${networkId}`)
      .then(res => parseNodesRes(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async deleteNodeFromOrg(args: NodeInput): Promise<DeleteNode> {
    return this.delete(`/${VERSION}/${NODES}/${args.id}/sites`)
      .then(() =>
        this.delete(`${args.id}`).then(() => {
          return { id: args.id };
        })
      )
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async attachNode(args: AttachNodeInput): Promise<CBooleanResponse> {
    return this.post(`/${VERSION}/${NODES}/${args.parentNode}/attach`, {
      body: {
        anodel: args.anodel,
        anoder: args.anoder,
      },
    })
      .then(res => (res ? { success: true } : { success: false }))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async detachhNode(args: NodeInput): Promise<CBooleanResponse> {
    return this.delete(`/${VERSION}/${NODES}/${args.id}/detach`)
      .then(res => (res ? { success: true } : { success: false }))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async addNode(args: AddNodeInput): Promise<Node> {
    return this.post(`/${VERSION}/${NODES}/`, {
      body: {
        name: args.name,
        node_id: args.id,
        org_id: args.orgId,
      },
    })
      .then(res => parseNodeRes(res.node))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async addNodeToSite(args: AddNodeToSiteInput): Promise<CBooleanResponse> {
    return this.post(`/${VERSION}/${NODES}/${args.nodeId}/sites`, {
      body: {
        net_id: args.networkId,
        site_id: args.siteId,
      },
    })
      .then(res => (res ? { success: true } : { success: false }))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async releaseNodeFromSite(args: NodeInput): Promise<CBooleanResponse> {
    return await this.delete(`/${VERSION}/${NODES}/${args.id}/sites`)
      .then(res => (res ? { success: true } : { success: false }))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async updateNodeState(args: UpdateNodeStateInput): Promise<Node> {
    return this.patch(`/${VERSION}/${NODES}/${args.id}/state/${args.state}`)
      .then(res => parseNodeRes(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
  async updateNode(args: UpdateNodeInput): Promise<Node> {
    return this.put(`/${VERSION}/${NODES}/${args.id}`, {
      body: {
        name: args.name,
      },
    })
      .then(res => parseNodeRes(res.node))
      .catch(err => {
        throw new GraphQLError(err);
      });
  }
}

export default NodeAPI;
