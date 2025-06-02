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
import { API_METHOD_TYPE, STATS_TYPE } from "../../common/enums";
import { logger } from "../../common/logger";
import {
  GetMetricsStatInput,
  MetricsRes,
  NotificationsRes,
} from "../resolvers/types";
import { parseMetricsResponse, parseNotificationsRes } from "./mapper";

const getNodeMetricRange = async (
  baseUrl: string,
  type: string,
  args: GetMetricsStatInput
): Promise<MetricsRes> => {
  const { to, from, nodeId, networkId = "", siteId } = args;
  let params = `from=${from}&to=${to}&step=1`;
  if (nodeId) {
    params = params + `&node=${nodeId}`;
  }
  if (networkId && type !== "subscribers_active") {
    params = params + `&network=${networkId}`;
  }
  if (siteId) {
    params = params + `&site=${siteId}`;
  }
  if (args.operation) {
    params = params + `&operation=${args.operation}`;
  }
  logger.info(
    `[getNodeMetricRange] Request URL: ${baseUrl}/${VERSION}/range/metrics/${type}?${params}`
  );
  return await asyncRestCall({
    method: API_METHOD_TYPE.GET,
    url: `${baseUrl}/${VERSION}/range/metrics/${type}?${params}`,
  })
    .then(res =>
      parseMetricsResponse(
        type === STATS_TYPE.DATA_USAGE ? res.data : res.data.data.result,
        type,
        args
      )
    )
    .catch(err => {
      logger.error(`Error fetching metrics: ${err}`);
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

export { getNodeMetricRange, getNotifications };
