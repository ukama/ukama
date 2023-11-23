/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

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
    })
      .then(res => dtoToNetworksDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getNetwork = async (networkId: string): Promise<NetworkDto> => {
    return this.get(`/${VERSION}/networks/${networkId}`)
      .then(res => dtoToNetworkDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSites = async (networkId: string): Promise<SitesResDto> => {
    return this.get(`/${VERSION}/networks/${networkId}/sites`)
      .then(res => dtoToSitesDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getSite = async (siteId: string, networkId: string): Promise<SiteDto> => {
    return this.get(`/${VERSION}/networks/${networkId}/sites/${siteId}`)
      .then(res => dtoToSiteDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  addNetwork = async (req: AddNetworkInputDto): Promise<NetworkDto> => {
    return this.post(`/${VERSION}/networks`, {
      body: req,
    })
      .then(res => dtoToNetworkDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  addSite = async (
    networkId: string,
    req: AddSiteInputDto
  ): Promise<SiteDto> => {
    return this.post(`/${VERSION}/networks/${networkId}/sites`, {
      body: req,
    })
      .then(res => dtoToSiteDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default NetworkApi;
