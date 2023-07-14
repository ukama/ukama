import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW } from "../../common/configs";

class NodeAPI extends RESTDataSource {
  baseURL = REGISTRY_API_GW;

  async getNode(args: GetNode) {
    return this.get(`nodes/${args.id}`);
  }
}

export default NodeAPI;
