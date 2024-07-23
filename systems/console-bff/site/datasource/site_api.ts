/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import { AddSiteInputDto, SiteDto, SitesResDto } from "../resolvers/types";
import { dtoToSiteDto, dtoToSitesDto } from "./mapper";

const SITES = "sites";

class SiteApi extends RESTDataSource {
  getSites = async (
    baseURL: string,
    networkId: string
  ): Promise<SitesResDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SITES}`, {
      params: {
        network: networkId,
      },
    }).then(res => dtoToSitesDto(res));
  };

  getSite = async (baseURL: string, siteId: string): Promise<SiteDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SITES}/${siteId}`).then(res =>
      dtoToSiteDto(res)
    );
  };

  addSite = async (baseURL: string, req: AddSiteInputDto): Promise<SiteDto> => {
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${SITES}`, {
      body: {
        access_id: req.access_id,
        backhaul_id: req.backhaul_id,
        install_date: req.install_date,
        is_deactivated: false,
        latitude: req.latitude,
        location: req.location,
        longitude: req.longitude,
        network_id: req.network_id,
        power_id: req.power_id,
        site: req.name,
        switch_id: req.switch_id,
        spectrum_id: req.spectrum_id,
      },
    }).then(res => dtoToSiteDto(res));
  };
}

export default SiteApi;
