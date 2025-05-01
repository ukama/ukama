/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { GraphQLError } from "graphql";

import { asyncRestCall } from "../../common/axiosClient";
import { VERSION } from "../../common/configs";
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
  parseNotificationsRes,
  parseSiteMetricRes,
} from "./mapper";

const getNodeMetricRange = async (
  baseUrl: string,
  type: string,
  args: GetMetricsStatInput | GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, step = 1, nodeId, userId } = args;
  let params = `from=${from}&step=${step}`;
  if (nodeId) {
    params = params + `&node=${nodeId}`;
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

const getSiteMetricRange = async (
  baseUrl: string,
  type: string,
  args: GetMetricRangeInput
): Promise<MetricRes> => {
  const { from, step = 1, userId, siteId } = args;

  let params = `from=${from}&step=${step}`;
  if (siteId) {
    params = params + `&site=${siteId}`;
  }
  if (userId) {
    params = params + `&user=${userId}`;
  }

  logger.info(
    `[getMetricRange] Request URL: ${baseUrl}/${VERSION}/range/metrics/${type}?${params}`
  );

  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${baseUrl}/${VERSION}/range/metrics/${type}?${params}`,
  })
    .then(res => {
      return parseSiteMetricRes(res, type, args);
    })
    .catch(err => {
      logger.error(
        `[getMetricRange] Error fetching metrics for ${type}: ${err}`
      );
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

export { getNodeMetricRange, getSiteMetricRange, getNotifications };
