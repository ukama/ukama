/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import { GetHealthReportInputDto, HealthInfo } from "../resolvers/types";
import { dtoToHealthInfo } from "./mapper";

const HEALTH = "health";

class HealthApi extends RESTDataSource {
  list = async (
    baseURL: string,
    req: GetHealthReportInputDto
  ): Promise<HealthInfo> => {
    this.baseURL = baseURL;
    this.logger.info(
      `GetHealthReport [GET]: ${baseURL}/${VERSION}/${HEALTH}/nodes/${req.nodeId}/list?filter=${req.filter}`
    );
    return this.get(
      `/${VERSION}/${HEALTH}/nodes/${req.nodeId}/list?filter=${req.filter}`
    ).then((res: any) => {
      return dtoToHealthInfo(res);
    });
  };
}

export default HealthApi;
