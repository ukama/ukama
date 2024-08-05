import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import {
  RestartNodeResDto,
  RestartNodesResDto,
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
  ): Promise<RestartNodesResDto> => {
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

  /**
   * Asynchronously restarts a site based on provided URLs and site ID.
   * @example
   * functionName('https://example.com', '12345')
   * { message: 'Site restarted successfully', status: 'success' }
   * @param {string} baseURL - The base URL of the site endpoint.
   * @param {string} siteId - The unique identifier of the site to restart.
   * @returns {Promise<RestartSiteResDto>} Object containing the status and message of the restart operation.
   * @description
   *   - Logs the full URL of the restart endpoint before making the request.
   *   - Updates the baseURL property of the object.
   */
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
