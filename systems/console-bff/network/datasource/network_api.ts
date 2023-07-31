import { RESTDataSource } from "@apollo/datasource-rest";

import setupLogger from "../../config/setupLogger";
import { SERVER } from "../../constants/endpoints";
import {
  AddNetworkInputDto,
  AddSiteInputDto,
  NetworkDto,
  NetworkStatusDto,
  NetworksResDto,
  SiteDto,
  SitesResDto,
} from "../types";
import {
  dtoToDto,
  dtoToNetworkDto,
  dtoToNetworksDto,
  dtoToSiteDto,
  dtoToSitesDto,
} from "./mapper";

const logger = setupLogger("service");

class NetworkApi extends RESTDataSource {
  getNetworkStatus = async (orgId: string): Promise<NetworkStatusDto> => {
    const resLiveNodes = await this.get(
      `${SERVER.ORG}/${orgId}/metrics/live-nodes`
    );
    const resTotalNodes = await this.get(`${SERVER.ORG}/${orgId}/nodes`);

    return dtoToDto(resTotalNodes.nodes.length, resLiveNodes.data.result[0]);
  };

  getNetworks = async (orgId: string): Promise<NetworksResDto> => {
    return this.get(`${SERVER.REGISTRY_NETWORKS_API_URL}`, {
      params: {
        org: orgId,
      },
    }).then(res => dtoToNetworksDto(res));
  };

  getNetwork = async (networkId: string): Promise<NetworkDto> => {
    return this.get(
      `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}`,
      {}
    ).then(res => dtoToNetworkDto(res));
  };

  getSites = async (networkId: string): Promise<SitesResDto> => {
    return this.get(
      `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites`,
      {}
    ).then(res => dtoToSitesDto(res));
  };

  getSite = async (siteId: string, networkId: string): Promise<SiteDto> => {
    return this.get(
      `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites/${siteId}`
    ).then(res => dtoToSiteDto(res));
  };

  addNetwork = async (req: AddNetworkInputDto): Promise<NetworkDto> => {
    return this.post(`${SERVER.REGISTRY_NETWORKS_API_URL}`, {
      body: req,
    }).then(res => dtoToNetworkDto(res));
  };

  addSite = async (
    networkId: string,
    req: AddSiteInputDto
  ): Promise<SiteDto> => {
    return this.post(`${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites`, {
      body: req,
    }).then(res => dtoToSiteDto(res));
  };
}

export default NetworkApi;
