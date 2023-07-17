import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW } from "../../common/configs";
import { BooleanResponse } from "../../common/types";

class NodeAPI extends RESTDataSource {
  baseURL = REGISTRY_API_GW;

  async getNode(args: NodeId): Promise<GetNode> {
    const node = this.get(`/v1/nodes/${args.id}`);
    return node;
  }
  async getNodes(): Promise<GetNodes> {
    const nodes = await this.get("/v1/nodes");
    return { nodes };
  }
  async getFreeNode(): Promise<GetNodes> {
    const nodes = await this.get("/v1/nodes/free");
    return { nodes };
  }
  async deleteNodeFromOrg(args: NodeId): Promise<DeleteNode> {
    const node = await this.delete(`/v1/nodes/${args.id}`);
    return { node };
  }
  async attachNode(args: AttachNodeArgs): Promise<BooleanResponse> {
    const node = await this.post("/v1/nodes/attach", {
      body: {
        anodel: args.anodel,
        anoder: args.anoder,
        parent_node: args.parentNode,
      },
    });
    return {
      success: node.success,
    };
  }
  async detachhNode(args: NodeId): Promise<BooleanResponse> {
    const node = await this.post("/v1/nodes/detach", {
      body: {
        node: args.id,
      },
    });
    return {
      success: node.success,
    };
  }
  async addNode(args: AddNodeArgs): Promise<Node> {
    const node = await this.post(`/v1/nodes/${args.id}`, {
      body: {
        state: args.state,
      },
    });
    return node;
  }
  async releaseNodeFromNetwork(args: NodeId): Promise<BooleanResponse> {
    const node = await this.post(`/v1/nodes/${args.id}/networks/release`);
    return { success: node.success };
  }
  async addNodeToNetwork(args: AddNodeToNetworkArgs): Promise<BooleanResponse> {
    const node = await this.post(
      `/v1/nodes/${args.nodeId}/networks/${args.networkId}/assign`
    );
    return { success: node.success };
  }
  async updateNodeState(args: UpdateNodeStateArgs): Promise<UpdateNodeState> {
    const node = await this.post(`/v1/nodes/${args.id}/state/${args.state}`);
    return { id: node.id, state: node.state };
  }
  async updateNode(args: UpdateNodeArgs): Promise<Node> {
    const node = await this.post(`/v1/nodes/${args.id}/update`, {
      body: {
        name: args.name,
      },
    });
    return node;
  }
}

export default NodeAPI;
