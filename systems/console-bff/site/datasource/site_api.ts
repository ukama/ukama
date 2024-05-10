/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW, VERSION } from "../../common/configs";
import { AddSiteInputDto, SiteDto, SitesResDto } from "../resolvers/types";
import { dtoToSiteDto, dtoToSitesDto } from "./mapper";

class SiteApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW;

  getSites = async (networkId: string): Promise<SitesResDto> => {
    return this.get(`/${VERSION}/networks/${networkId}/sites`).then(res =>
      dtoToSitesDto(res)
    );
  };

  getSite = async (siteId: string, networkId: string): Promise<SiteDto> => {
    return this.get(`/${VERSION}/networks/${networkId}/sites/${siteId}`).then(
      res => dtoToSiteDto(res)
    );
  };

  addSite = async (
    networkId: string,
    req: AddSiteInputDto
  ): Promise<SiteDto> => {
    return this.post(`/${VERSION}/networks/${networkId}/sites`, {
      body: req,
    }).then(res => dtoToSiteDto(res));
  };
}

export default SiteApi;
