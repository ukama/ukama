/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  NetworkAPIResDto,
  NetworkDto,
  NetworksAPIResDto,
  NetworksResDto,
} from "../resolvers/types";

export const dtoToNetworksDto = (res: NetworksAPIResDto): NetworksResDto => {
  const networks: NetworkDto[] = [];
  res.networks.forEach(network => {
    networks.push({
      id: network.id,
      name: network.name,
      budget: network.budget,
      isDefault: false,
      isDeactivated: network.is_deactivated,
      createdAt: network.created_at,
      countries: network.allowed_countries,
      networks: network.allowed_networks,
      overdraft: network.overdraft,
      paymentLinks: network.payment_links,
      trafficPolicy: network.traffic_policy,
    });
  });
  return {
    networks: networks,
  };
};

export const dtoToNetworkDto = (res: NetworkAPIResDto): NetworkDto => {
  return {
    id: res.network.id,
    name: res.network.name,
    budget: res.network.budget,
    isDefault: false,
    isDeactivated: res.network.is_deactivated,
    createdAt: res.network.created_at,
    countries: res.network.allowed_countries,
    networks: res.network.allowed_networks,
    overdraft: res.network.overdraft,
    paymentLinks: res.network.payment_links,
    trafficPolicy: res.network.traffic_policy,
  };
};
