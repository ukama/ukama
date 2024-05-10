/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW, VERSION } from "../../common/configs";
import {
  AddNetworkInputDto,
  NetworkDto,
  NetworksResDto,
} from "../resolvers/types";
import { dtoToNetworkDto, dtoToNetworksDto } from "./mapper";

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

  addNetwork = async (req: AddNetworkInputDto): Promise<NetworkDto> => {
    return this.post(`/${VERSION}/networks`, {
      body: {
        allowed_countries: req.countries,
        allowed_networks: req.networks,
        budget: req.budget,
        network_name: req.name,
        org: req.org,
        overdraft: 0,
        payment_links: true,
        traffic_policy: 0,
      },
    }).then(res => dtoToNetworkDto(res));
  };
}

export default NetworkApi;
