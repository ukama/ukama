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
import {
  METRIC_API_GW,
  METRIC_PROMETHEUS,
  NOTIFICATION_API_GW,
} from "../../common/configs";
import { API_METHOD_TYPE } from "../../common/enums";
import {
  GetLatestMetricInput,
  GetMetricRangeInput,
  GetNotificationsInput,
  LatestMetricRes,
  MetricRes,
  NotificationsRes,
} from "../resolvers/types";
import {
  parseLatestMetricRes,
  parseMetricRes,
  parseNodeMetricRes,
  parseNotificationsRes,
  parsePromethRes,
} from "./mapper";

const getLatestMetric = async (
  args: GetLatestMetricInput
): Promise<LatestMetricRes> => {
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${METRIC_API_GW}/v1/metrics/${args.type}`,
  }).then(res => parseLatestMetricRes(res.data, args));
};

const directCall = async (args: GetMetricRangeInput): Promise<MetricRes> => {
  const { from, to, step = 1 } = args;
  const agent = new https.Agent({
    rejectUnauthorized: false,
  });
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    httpsAgent: agent,
    url: `${METRIC_PROMETHEUS}?query=${args.type}&start=${from}&end=${to}&step=${step}`,
  })
    .then(res => parsePromethRes(res.data, args))
    .catch(err => {
      throw new GraphQLError(err);
    });
};

const getMetricRange = async (
  args: GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, to = 0, step = 1 } = args;
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${METRIC_API_GW}/v1/range/metrics/${args.type}?from=${from}&to=${to}&step=${step}`,
  })
    .then(res => parseMetricRes(res.data, args.type))
    .catch(err => {
      throw new GraphQLError(err);
    });
};

const getNodeRangeMetric = async (
  args: GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, to = 0, step = 1 } = args;
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${METRIC_API_GW}/v1/nodes/${args.nodeId}/metrics/${args.type}?from=${from}&to=${to}&step=${step}`,
  })
    .then(res => parseNodeMetricRes(res, args))
    .catch(err => {
      throw new GraphQLError(err);
    });
};

const getNotifications = async (
  args: GetNotificationsInput
): Promise<NotificationsRes> => {
  const { orgId, subscriberId, userId, networkId, forRole } = args;

  let params = "";
  if (orgId) {
    params = params + `&org_id=${orgId}`;
  }
  if (subscriberId) {
    params = params + `&subscriber_id=${subscriberId}`;
  }
  if (userId) {
    params = params + `&user_id=${userId}`;
  }
  if (networkId) {
    params = params + `&network_id=${networkId}`;
  }
  if (forRole) {
    params = params + `&role=${forRole}`;
  }
  if (params.length > 0) params = params.substring(1);

  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${NOTIFICATION_API_GW}/v1/event-notification?${params}`,
  }).then(res => parseNotificationsRes(res.data));
};

export {
  directCall,
  getLatestMetric,
  getMetricRange,
  getNodeRangeMetric,
  getNotifications,
};
