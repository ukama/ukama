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
import { VERSION } from "../../common/configs";
import { API_METHOD_TYPE } from "../../common/enums";
import { logger } from "../../common/logger";
import {
  GetMetricRangeInput,
  GetNotificationsInput,
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
  args: GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, to = 0, step = 1 } = args;
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${baseUrl}/${VERSION}/range/metrics/${args.type}?from=${from}&to=${to}&step=${step}`,
  })
    .then(res => parseMetricRes(res.data, args.type))
    .catch(err => {
      throw new GraphQLError(err);
    });
};

const getNodeRangeMetric = async (
  baseUrl: string,
  args: GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, to = 0, step = 1 } = args;
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
  args: GetNotificationsInput
): Promise<NotificationsRes> => {
  const { orgId, userId, networkId, subscriberId, nodeId } = args;

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

  if (subscriberId) {
    params = params + `&node_id=${nodeId}`;
  }
  if (params.length > 0) params = params.substring(1);
  logger.info(
    `GetNotifications [GET]: ${baseUrl}/${VERSION}/event-notification?${params}`
  );
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${baseUrl}/${VERSION}/event-notification?${params}`,
  }).then(res => {
    return parseNotificationsRes(res.data);
  });
};

export { directCall, getMetricRange, getNodeRangeMetric, getNotifications };
