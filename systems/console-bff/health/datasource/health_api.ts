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
    const query = new URLSearchParams();
    query.set("timeframe", req.timeframe || "all");
    if (req.id) query.set("id", req.id);
    if (req.nodeId) query.set("node_id", req.nodeId);
    if (req.timestamp) query.set("timestamp", req.timestamp);
    this.logger.info(
      `GetHealthReport [GET]: ${baseURL}/${VERSION}/${HEALTH}/list?${query.toString()}`
    );
    return this.get(`/${VERSION}/${HEALTH}/list?${query.toString()}`).then(
      (res: any) => {
        return dtoToHealthInfo(res);
      }
    );
  };
}

export default HealthApi;
