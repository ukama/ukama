/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { SITE_ADD_TIMEOUT_MS, VERSION } from "../../common/configs";
import { TimedRESTDataSource } from "../../common/datasource";
import { AddSiteInputDto, SiteDto } from "../resolvers/types";
import { dtoToSiteDto } from "./mapper";

const SITES = "sites";

class SiteProvisionApi extends TimedRESTDataSource {
  constructor() {
    super(SITE_ADD_TIMEOUT_MS);
  }

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
}

export default SiteProvisionApi;
