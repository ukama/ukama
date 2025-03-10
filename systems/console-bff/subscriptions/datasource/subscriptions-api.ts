/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { GraphQLError } from "graphql";
import https from "https";

import { asyncRestCall } from "../../common/axiosClient";
import { METRICS_INTERVAL, VERSION } from "../../common/configs";
import { API_METHOD_TYPE } from "../../common/enums";
import { logger } from "../../common/logger";
import {
  GetMetricRangeInput,
  GetMetricsStatInput,
  MetricRes,
  NotificationsRes,
} from "../resolvers/types";
import {
  parseMetricRes,
  parseNodeMetricRes,
  parseNotificationsRes,
  parsePromethRes,
} from "./mapper";

const directCall = async (
  baseUrl: string,
  args: GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, to, step = 1 } = args;
  const agent = new https.Agent({
    rejectUnauthorized: false,
  });
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    httpsAgent: agent,
    url: `${baseUrl}?query=${args.type}&start=${from}&end=${to}&step=${step}`,
  })
    .then(res => parsePromethRes(res.data, args))
    .catch(err => {
      throw new GraphQLError(err);
    });
};

const getMetricRange = async (
  baseUrl: string,
  type: string,
  args: GetMetricsStatInput
): Promise<MetricRes> => {
  const { from, step = 1, nodeId, userId, networkId } = args;
  let params = `from=${from}&step=${step}`;
  if (nodeId) {
    params = params + `&node=${nodeId}`;
  }
  if (networkId) {
    params = params + `&network=${networkId}`;
  }
  if (userId) {
    params = params + `&user=${userId}`;
  }
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${baseUrl}/${VERSION}/range/metrics/${type}?${params}`,
  })
    .then(res => parseMetricRes(res, type, args))
    .catch(err => {
      throw new GraphQLError(err);
    });
};

const getNodeRangeMetric = async (
  baseUrl: string,
  args: GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, to = 0, step = METRICS_INTERVAL } = args;
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${baseUrl}/${VERSION}/nodes/${args.nodeId}/metrics/${args.type}?from=${from}&to=${to}&step=${step}`,
  })
    .then(res => parseNodeMetricRes(res, args))
    .catch(err => {
      throw new GraphQLError(err);
    });
};

const getNotifications = async (
  baseUrl: string,
  orgId: string,
  userId: string,
  networkId: string,
  subscriberId: string
): Promise<NotificationsRes> => {
  let params = "";
  if (orgId) {
    params = params + `&org_id=${orgId}`;
  }
  if (userId) {
    params = params + `&user_id=${userId}`;
  }
  if (networkId) {
    params = params + `&network_id=${networkId}`;
  }
  if (subscriberId) {
    params = params + `&subscriber_id=${subscriberId}`;
  }

  if (params.length > 0) params = params.substring(1);
  logger.info(
    `GetNotifications [GET]: ${baseUrl}/${VERSION}/event-notification?${params}`
  );
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${baseUrl}/${VERSION}/event-notification?${params}`,
  }).then(res => parseNotificationsRes(res.data));
};

export { directCall, getMetricRange, getNodeRangeMetric, getNotifications };
