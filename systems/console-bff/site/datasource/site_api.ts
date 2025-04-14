/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import {
  AddSiteInputDto,
  SiteDto,
  SitesInputDto,
  SitesResDto,
  UpdateSiteInputDto,
} from "../resolvers/types";
import { dtoToSiteDto, dtoToSitesDto } from "./mapper";

const SITES = "sites";

class SiteApi extends RESTDataSource {
  getSites = async (
    baseURL: string,
    args: SitesInputDto
  ): Promise<SitesResDto> => {
    let params = "";
    if (args.networkId) {
      params = params.concat(`network_id=${args.networkId}`);
    }
    this.logger.info(
      `GetSites [GET]: ${baseURL}/${VERSION}/${SITES}?${params}`
    );
    this.baseURL = baseURL;

    return this.get(`/${VERSION}/${SITES}?${params}`).then(res =>
      dtoToSitesDto(res)
    );
  };

  getSite = async (baseURL: string, siteId: string): Promise<SiteDto> => {
    this.logger.info(`GetSite [GET]: ${baseURL}/${VERSION}/${SITES}/${siteId}`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SITES}/${siteId}`).then(res =>
      dtoToSiteDto(res)
    );
  };

  addSite = async (baseURL: string, req: AddSiteInputDto): Promise<SiteDto> => {
    this.logger.info(`AddSite [POST]: ${baseURL}/${VERSION}/${SITES}`);
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

  updateSite = async (
    baseURL: string,
    siteId: string,
    req: UpdateSiteInputDto
  ): Promise<SiteDto> => {
    this.logger.info(
      `UpdateSite [PATCH]: ${baseURL}/${VERSION}/${SITES}/${siteId}`
    );
    this.baseURL = baseURL;
    return this.patch(`/${VERSION}/${SITES}/${siteId}`, {
      body: {
        name: req.name,
      },
    }).then(res => dtoToSiteDto(res));
  };
}

export default SiteApi;
