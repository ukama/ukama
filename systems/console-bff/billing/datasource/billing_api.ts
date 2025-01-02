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
  GetReportDto,
  GetReportsDto,
  GetReportsInputDto,
} from "../resolvers/types";
import { dtoToReportDto, dtoToReportsDto } from "./mapper";

class BillingAPI extends RESTDataSource {
  getReports = async (
    baseURL: string,
    req: GetReportsInputDto
  ): Promise<GetReportsDto> => {
    this.logger.info(`GetReports [GET] ${baseURL}/${VERSION}/reports`);
    this.baseURL = baseURL;
    const params = "sort=true";
    if (req.networkId) {
      params.concat(`&network_id=${req.networkId}`);
    }
    if (req.ownerId) {
      params.concat(`&owner_id=${req.ownerId}`);
    }
    if (req.ownerType) {
      params.concat(`&owner_type=${req.ownerType}`);
    }
    if (req.report_type) {
      params.concat(`&report_type=${req.report_type}`);
    }
    if (req.count) {
      params.concat(`&count=${req.count}`);
    }
    if (req.isPaid) {
      params.concat(`&is_paid=${req.isPaid}`);
    }
    return this.get(`/${VERSION}/reports?${params}`).then(res =>
      dtoToReportsDto(res)
    );
  };

  getReport = async (baseURL: string, id: string): Promise<GetReportDto> => {
    this.logger.info(`GetReport [GET]: ${baseURL}/${VERSION}/report/${id}`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/reports/${id}`).then(res =>
      dtoToReportDto(res)
    );
  };
  getPDFReport = async (baseURL: string, id: string): Promise<GetReportDto> => {
    this.logger.info(`GetReport [GET]: ${baseURL}/${VERSION}/reports/${id}`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/pdf/${id}?as_pdf=true`).then(res => res);
  };
}

export default BillingAPI;
