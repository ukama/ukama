import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW } from "../../common/configs";
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

const version = "/v1/networks";

class NetworkApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW + version;
  getNetworkStatus = async (orgId: string): Promise<NetworkStatusDto> => {
    const resLiveNodes = await this.get(`/metrics/live-nodes`);
    const resTotalNodes = await this.get(`${SERVER.ORG}/${orgId}/nodes`);

    return dtoToDto(resTotalNodes.nodes.length, resLiveNodes.data.result[0]);
  };

  getNetworks = async (orgId: string): Promise<NetworksResDto> => {
    return this.get("", {
      params: {
        org: orgId,
      },
    }).then(res => dtoToNetworksDto(res));
  };

  getNetwork = async (networkId: string): Promise<NetworkDto> => {
    return this.get(`/${networkId}`).then(res => dtoToNetworkDto(res));
  };

  getSites = async (networkId: string): Promise<SitesResDto> => {
    return this.get(`/${networkId}/sites`).then(res => dtoToSitesDto(res));
  };

  getSite = async (siteId: string, networkId: string): Promise<SiteDto> => {
    return this.get(`/${networkId}/sites/${siteId}`).then(res =>
      dtoToSiteDto(res)
    );
  };

  addNetwork = async (req: AddNetworkInputDto): Promise<NetworkDto> => {
    return this.post("", {
      body: req,
    }).then(res => dtoToNetworkDto(res));
  };

  addSite = async (
    networkId: string,
    req: AddSiteInputDto
  ): Promise<SiteDto> => {
    return this.post(`/${networkId}/sites`, {
      body: req,
    }).then(res => dtoToSiteDto(res));
  };
}

export default NetworkApi;
