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
      isDeactivated: network.is_deactivated,
      createdAt: network.created_at,
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
    isDeactivated: res.network.is_deactivated,
    createdAt: res.network.created_at,
  };
};
