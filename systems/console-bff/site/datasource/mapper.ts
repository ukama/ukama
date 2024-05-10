/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
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
