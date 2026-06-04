/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { BaseRESTDataSource } from "../../common/datasource";
import {
  GetNodeLatestMetricInput,
  GetSiteLatestMetricInput,
  NodeLatestMetric,
  SiteLatestMetric,
} from "../resolver/types";
import { parseNodeLatestMetricRes, parseSiteLatestMetricRes } from "./mapper";

const VERSION = "v1";
const METRICS = "metrics";

class MetricAPI extends BaseRESTDataSource {
  getNodeLatestMetric = async (
    baseURL: string,
    args: GetNodeLatestMetricInput
  ): Promise<NodeLatestMetric> => {
    this.logger.info(
      `GetNodeLatestMetric [GET]: ${baseURL}/${VERSION}/${METRICS}/${args.type}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${METRICS}/${args.type}`).then(res =>
      parseNodeLatestMetricRes(res, args)
    );
  };

  getSiteLatestMetric = async (
    baseURL: string,
    args: GetSiteLatestMetricInput
  ): Promise<SiteLatestMetric> => {
    this.logger.info(
      `GetSiteLatestMetric [GET]: ${baseURL}/${VERSION}/${METRICS}/${args.type}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${METRICS}/${args.type}`).then(res =>
      parseSiteLatestMetricRes(res, args)
    );
  };

  /**
   * Generic latest-value read for one metric key (org-scoped, no entity
   * stamping). Used by the dashboard KPI sections (plan Phase 4 — polled,
   * no subscriptions in v1).
   */
  getLatestMetric = async (
    baseURL: string,
    type: string
  ): Promise<{ type: string; value: [number, number]; success: boolean }> => {
    this.logger.info(
      `GetLatestMetric [GET]: ${baseURL}/${VERSION}/${METRICS}/${type}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${METRICS}/${type}`).then(res => {
      const data = res?.data?.result?.[0];
      if (data?.value?.length > 0) {
        return { type, value: data.value as [number, number], success: true };
      }
      return { type, value: [0, 0] as [number, number], success: false };
    });
  };
}

export default MetricAPI;
