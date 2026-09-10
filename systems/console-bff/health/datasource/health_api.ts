/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { VERSION } from "../../common/configs";
import { BaseRESTDataSource } from "../../common/datasource";
import { TIMEFRAME_FILTER } from "../../common/enums";
import {
  Apps,
  GetAppsInputDto,
  GetHealthReportInputDto,
  GetNodeInterfacesInputDto,
  HealthInfo,
  NodeInterfaces,
} from "../resolvers/types";
import { mapApps, mapNodeInterfaces, mapReportsToHealthInfo } from "./mapper";

const HEALTH = "health";

const nodePath = (nodeId: string) =>
  `/${VERSION}/${HEALTH}/nodes/${encodeURIComponent(nodeId)}`;

class HealthApi extends BaseRESTDataSource {
  // Per-app runtime health/resource list for a node, served by the node
  // api-gateway. appName narrows to a single app.
  getApps = async (baseURL: string, data: GetAppsInputDto): Promise<Apps> => {
    const { nodeId, appName } = data;
    const queryParams = new URLSearchParams();
    if (appName) {
      queryParams.append("appName", appName);
    }
    const query = queryParams.toString();
    const path = `${nodePath(nodeId)}/apps${query ? `?${query}` : ""}`;
    this.baseURL = baseURL;
    this.logger.info(`GetApps [GET]: ${baseURL}${path}`);
    return this.get(path)
      .then(apps => mapApps(apps))
      .catch(error => {
        this.logger.error(`Error getting apps: ${error}`);
        throw error;
      });
  };

  getInterfaces = async (
    baseURL: string,
    data: GetNodeInterfacesInputDto
  ): Promise<NodeInterfaces> => {
    const path = `${nodePath(data.nodeId)}/interfaces`;
    this.baseURL = baseURL;
    this.logger.info(`GetNodeInterfaces [GET]: ${baseURL}${path}`);
    return this.get(path)
      .then(res => mapNodeInterfaces(data.nodeId, res))
      .catch(error => {
        this.logger.error(`Error getting node interfaces: ${error}`);
        throw error;
      });
  };

  list = async (
    baseURL: string,
    req: GetHealthReportInputDto
  ): Promise<HealthInfo> => {
    this.baseURL = baseURL;
    const query = new URLSearchParams();
    query.set(
      "timeframe",
      req.timeframe === TIMEFRAME_FILTER.LATEST
        ? TIMEFRAME_FILTER.LATEST
        : TIMEFRAME_FILTER.ALL
    );
    if (req.id) query.set("reportId", req.id);
    if (req.timestamp) query.set("reportedAt", req.timestamp);
    const path = `${nodePath(req.nodeId)}/reports?${query.toString()}`;
    this.logger.info(`GetHealthReport [GET]: ${baseURL}${path}`);
    return this.get(path).then(res => mapReportsToHealthInfo(res));
  };
}

export default HealthApi;
