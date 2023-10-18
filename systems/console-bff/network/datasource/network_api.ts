import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW, VERSION } from "../../common/configs";
import {
  AddNetworkInputDto,
  AddSiteInputDto,
  NetworkDto,
  NetworksResDto,
  SiteDto,
  SitesResDto,
} from "../resolvers/types";
import {
  dtoToNetworkDto,
  dtoToNetworksDto,
  dtoToSiteDto,
  dtoToSitesDto,
} from "./mapper";

class NetworkApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW;

  getNetworks = async (orgId: string): Promise<NetworksResDto> => {
    return this.get(`/${VERSION}/networks`, {
      params: {
        org: orgId,
      },
    }).then(res => dtoToNetworksDto(res));
  };

  getNetwork = async (networkId: string): Promise<NetworkDto> => {
    return this.get(`/${VERSION}/networks/${networkId}`).then(res =>
      dtoToNetworkDto(res)
    );
  };

  getSites = async (networkId: string): Promise<SitesResDto> => {
    return this.get(`/${VERSION}/networks/${networkId}/sites`).then(res =>
      dtoToSitesDto(res)
    );
  };

  getSite = async (siteId: string, networkId: string): Promise<SiteDto> => {
    return this.get(`/${VERSION}/networks/${networkId}/sites/${siteId}`).then(
      res => dtoToSiteDto(res)
    );
  };

  addNetwork = async (req: AddNetworkInputDto): Promise<NetworkDto> => {
    return this.post(`/${VERSION}/networks`, {
      body: req,
    }).then(res => dtoToNetworkDto(res));
  };

  addSite = async (
    networkId: string,
    req: AddSiteInputDto
  ): Promise<SiteDto> => {
    return this.post(`/${VERSION}/networks/${networkId}/sites`, {
      body: req,
    }).then(res => dtoToSiteDto(res));
  };
}

export default NetworkApi;
