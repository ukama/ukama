/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { VERSION } from "../../common/configs";
import { BaseRESTDataSource } from "../../common/datasource";
import {
  GetReportDto,
  GetReportsDto,
  GetReportsInputDto,
} from "../resolvers/types";
import { dtoToReportDto, dtoToReportsDto } from "./mapper";

class BillingAPI extends BaseRESTDataSource {
  getReports = async (
    baseURL: string,
    req: GetReportsInputDto
  ): Promise<GetReportsDto> => {
    this.logger.info(`GetReports [GET] ${baseURL}/${VERSION}/reports`);
    this.baseURL = baseURL;
    const params = new URLSearchParams({ sort: "true" });
    if (req.networkId) {
      params.append("network_id", req.networkId);
    }
    if (req.ownerId) {
      params.append("owner_id", req.ownerId);
    }
    if (req.ownerType) {
      params.append("owner_type", req.ownerType);
    }
    if (req.report_type) {
      params.append("report_type", req.report_type);
    }
    if (req.count) {
      params.append("count", `${req.count}`);
    }
    if (req.isPaid !== undefined && req.isPaid !== null) {
      params.append("is_paid", `${req.isPaid}`);
    }
    return this.get(`/${VERSION}/reports?${params.toString()}`).then(res =>
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
