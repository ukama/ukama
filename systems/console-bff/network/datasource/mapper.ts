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
  SiteAPIResDto,
  SiteDto,
  SitesAPIResDto,
  SitesResDto,
} from "../resolvers/types";

export const dtoToSitesDto = (res: SitesAPIResDto): SitesResDto => {
  const sites: SiteDto[] = [];
  res.sites.forEach(site => {
    sites.push({
      id: site.id,
      name: site.name,
      networkId: site.network_id,
      isDeactivated: site.is_deactivated,
      createdAt: site.created_at,
    });
  });
  return {
    networkId: res.network_id,
    sites: sites,
  };
};

export const dtoToSiteDto = (res: SiteAPIResDto): SiteDto => {
  return {
    id: res.site.id,
    name: res.site.name,
    networkId: res.site.network_id,
    isDeactivated: res.site.is_deactivated,
    createdAt: res.site.created_at,
  };
};

export const dtoToNetworksDto = (res: NetworksAPIResDto): NetworksResDto => {
  const networks: NetworkDto[] = [];
  res.networks.forEach(network => {
    networks.push({
      id: network.id,
      name: network.name,
      orgId: network.org_id,
      budget: network.budget,
      isDeactivated: network.is_deactivated,
      createdAt: network.created_at,
      countries: network.allowed_countries,
      networks: network.allowed_networks,
    });
  });
  return {
    orgId: res.org_id,
    networks: networks,
  };
};

export const dtoToNetworkDto = (res: NetworkAPIResDto): NetworkDto => {
  return {
    id: res.network.id,
    name: res.network.name,
    orgId: res.network.org_id,
    budget: res.network.budget,
    isDeactivated: res.network.is_deactivated,
    createdAt: res.network.created_at,
    countries: res.network.allowed_countries,
    networks: res.network.allowed_networks,
  };
};
