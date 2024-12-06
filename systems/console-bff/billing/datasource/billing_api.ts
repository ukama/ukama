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
  GetReportInputDto,
  GetReportResDto,
  GetReportsInputDto,
  GetReportsResDto,
  InvoiceInputDto,
} from "../resolvers/types";
import { dtoToReportDto, dtoToReportsDto } from "./mapper";

class BillingAPI extends RESTDataSource {
  getReports = async (
    baseURL: string,
    req: GetReportsInputDto
  ): Promise<GetReportsResDto> => {
    let params = "";

    if (req.count) {
      params += `&count=${req.count}`;
    }
    if (req.is_paid !== undefined) {
      params += `&is_paid=${req.is_paid}`;
    }
    if (req.network_id) {
      params += `&network_id=${req.network_id}`;
    }
    if (req.owner_id) {
      params += `&owner_id=${req.owner_id}`;
    }
    if (req.owner_type) {
      params += `&owner_type=${req.owner_type}`;
    }
    if (req.report_type) {
      params += `&report_type=${req.report_type}`;
    }

    if (params.length > 0) {
      params = params.substring(1);
    }

    this.logger.info(
      `GetReports [GET]: ${baseURL}/${VERSION}/reports?${params}`
    );
    this.baseURL = baseURL;

    return this.get(`/${VERSION}/reports?${params}`).then(res =>
      dtoToReportsDto(res)
    );
  };

  getReport = async (
    baseURL: string,
    req: GetReportInputDto
  ): Promise<GetReportResDto> => {
    this.logger.info(`GetReport [GET]: ${baseURL}/${VERSION}/report/${req.id}`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/reports/${req.id}?as_pdf=${req.asPdf}`).then(
      res => dtoToReportDto(res)
    );
  };

  addReport = async (
    baseURL: string,
    req: InvoiceInputDto
  ): Promise<GetReportResDto> => {
    this.logger.info(`AddReport [POST]: ${baseURL}/${VERSION}/report`);
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/reports`, { body: req }).then(res => res);
  };
}

export default BillingAPI;
