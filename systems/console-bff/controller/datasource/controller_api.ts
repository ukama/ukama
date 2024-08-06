import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import {
  EmptyResDto,
  RestartNodeResDto,
  RestartSiteResDto,
} from "../resolvers/types";

const CONTROLLER = "controller";

class ControllerApi extends RESTDataSource {
  restartNode = async (
    baseURL: string,
    nodeId: string
  ): Promise<RestartNodeResDto> => {
    this.logger.info(
      `RestartNode [POST]: ${baseURL}/${VERSION}/${CONTROLLER}/${nodeId}/restart`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${CONTROLLER}/nodes/${nodeId}/restart`);
  };

  restartNodes = async (
    baseURL: string,
    nodeIds: string[]
  ): Promise<EmptyResDto> => {
    this.logger.info(
      `RestartNodes [POST]: ${baseURL}/${VERSION}/${CONTROLLER}/nodes/restart`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${CONTROLLER}/nodes/restart`, {
      body: {
        node_ids: nodeIds,
      },
    });
  };

  restartSite = async (
    baseURL: string,
    siteId: string
  ): Promise<RestartSiteResDto> => {
    this.logger.info(
      `${baseURL}/${VERSION}/${CONTROLLER}/sites/${siteId}/restart`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${CONTROLLER}/sites/${siteId}/restart`);
  };
}

export default ControllerApi;
