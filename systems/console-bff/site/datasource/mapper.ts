import {
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
      latitude: site.latitude,
      longitude: site.longitude,
      installDate: site.install_date,
      accessId: site.access_id,
      powerId: site.power_id,
      switchId: site.switch_id,
      backhaulId: site.backhaul_id,
      createdAt: site.created_at,
      location: site.location,
      spectrumId: site.spectrum_id,
    });
  });
  return {
    sites: sites,
  };
};

export const dtoToSiteDto = (res: SiteAPIResDto): SiteDto => {
  return {
    id: res.site.id,
    name: res.site.name,
    networkId: res.site.network_id,
    isDeactivated: res.site.is_deactivated,
    backhaulId: res.site.backhaul_id,
    switchId: res.site.switch_id,
    powerId: res.site.power_id,
    longitude: res.site.longitude,
    latitude: res.site.latitude,
    accessId: res.site.access_id,
    installDate: res.site.install_date,
    createdAt: res.site.created_at,
    location: res.site.location,
    spectrumId: res.site.spectrum_id,
  };
};
