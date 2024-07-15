import {
  SiteAPIDto,
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
      latitude: site.latitude,
      longitude: site.longitude,
      installDate: site.install_date,
      accessId: site.access_id,
      powerId: site.power_id,
      switchId: site.switch_id,
      backhaulId: site.backhaul_id,
      createdAt: site.created_at,
      location: site.location,
    });
  });
  return {
    sites: sites,
  };
};

export const dtoToSiteDto = (res: SiteAPIDto): SiteDto => {
  return {
    id: res.id,
    name: res.name,
    networkId: res.network_id,
    isDeactivated: res.is_deactivated,
    backhaulId: res.backhaul_id,
    switchId: res.switch_id,
    powerId: res.power_id,
    longitude: res.longitude,
    latitude: res.latitude,
    accessId: res.access_id,
    installDate: res.install_date,
    createdAt: res.created_at,
    location: res.location,
  };
};
