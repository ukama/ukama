/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { GetNodeLatestMetricInput, NodeLatestMetric } from "../resolver/types";
import { parseNodeLatestMetricRes } from "./mapper";

const VERSION = "v1";
const METRICS = "metrics";

class MetricAPI extends RESTDataSource {
  getNodeLatestMetric = async (
    baseURL: string,
    args: GetNodeLatestMetricInput
  ): Promise<NodeLatestMetric> => {
    this.logger.info(`GET: ${baseURL}/${VERSION}/${METRICS}/${args.type}`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${METRICS}/${args.type}`).then(res =>
      parseNodeLatestMetricRes(res, args)
    );
  };
}

export default MetricAPI;
