/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { VERSION } from "../../common/configs";
import { BaseRESTDataSource } from "../../common/datasource";
import {
  Apps,
  GetAppsInputDto,
  GetHealthReportInputDto,
  HealthInfo,
} from "../resolvers/types";
import { dtoToHealthInfo, mapApps } from "./mapper";

const HEALTH = "health";

class HealthApi extends BaseRESTDataSource {
  // Per-app runtime health/resource list for a node. nodeId is required;
  // appName narrows to a single app. baseURL resolves to nodeGwIp:nodeGwPort
  // (health is isForNodeGw) — the node gateway serves /v1/health/apps.
  getApps = async (baseURL: string, data: GetAppsInputDto): Promise<Apps> => {
    const { nodeId, appName } = data;
    const queryParams = new URLSearchParams();
    queryParams.append("nodeId", nodeId);
    if (appName) {
      queryParams.append("appName", appName);
    }
    this.baseURL = baseURL;
    this.logger.info(
      `GetApps [GET]: ${baseURL}/${VERSION}/${HEALTH}/apps?${queryParams.toString()}`
    );
    return this.get(`/${VERSION}/${HEALTH}/apps?${queryParams.toString()}`)
      .then(apps => mapApps(apps))
      .catch(error => {
        this.logger.error(`Error getting apps: ${error}`);
        throw error;
      });
  };

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
