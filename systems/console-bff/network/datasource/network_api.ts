/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import { TBooleanResponse } from "../../common/types";
import {
  AddNetworkInputDto,
  NetworkDto,
  NetworksResDto,
} from "../resolvers/types";
import { dtoToNetworkDto, dtoToNetworksDto } from "./mapper";

class NetworkApi extends RESTDataSource {
  getNetworks = async (baseURL: string): Promise<NetworksResDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/networks`).then(res => dtoToNetworksDto(res));
  };

  getNetwork = async (
    baseURL: string,
    networkId: string
  ): Promise<NetworkDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/networks/${networkId}`).then(res =>
      dtoToNetworkDto(res)
    );
  };

  addNetwork = async (
    baseURL: string,
    req: AddNetworkInputDto
  ): Promise<NetworkDto> => {
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/networks`, {
      body: {
        allowed_countries: req.countries,
        allowed_networks: req.networks,
        budget: req.budget,
        network_name: req.name,
        overdraft: 0,
        payment_links: true,
        traffic_policy: 0,
      },
    }).then(res => dtoToNetworkDto(res));
  };

  setDefaultNetwork = async (
    baseURL: string,
    networkId: string
  ): Promise<TBooleanResponse> => {
    this.baseURL = baseURL;
    return this.patch(`/${VERSION}/networks/${networkId}`)
      .then(() => {
        return { success: true };
      })
      .catch(() => {
        return { success: false };
      });
  };
}

export default NetworkApi;
