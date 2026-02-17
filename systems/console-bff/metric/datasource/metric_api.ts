/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { MetricAnalysis, MetricDomain } from "../resolver/types";
import { parseMetricAnalysisRes, parseMetricDomainRes } from "./mapper";

const VERSION = "v1";

class MetricAPI extends RESTDataSource {
  async getMetricAnalysis(
    baseURL: string,
    metricId: string,
    nodeId: string
  ): Promise<MetricAnalysis> {
    this.logger.info(
      `GetMetricAnalysis [GET]: ${baseURL}/${VERSION}/reasoning/stats/nodes/${nodeId}/metrics/${metricId}`
    );
    this.baseURL = "http://api-gateway-metrics:8080";
    return this.get(
      `/${VERSION}/reasoning/stats/nodes/${nodeId}/metrics/${metricId}`
    ).then(res => parseMetricAnalysisRes(res));
  }

  async getMetricDomain(
    baseURL: string,
    metricId: string,
    nodeId: string
  ): Promise<MetricDomain> {
    this.logger.info(
      `GetMetricDomain [GET]: ${baseURL}/${VERSION}/reasoning/domain/nodes/${nodeId}/metrics/${metricId}`
    );
    this.baseURL = "http://api-gateway-metrics:8080";
    return this.get(
      `/${VERSION}/reasoning/domain/nodes/${nodeId}/metrics/${metricId}`
    ).then(res => parseMetricDomainRes(res));
  }
}

export default MetricAPI;
