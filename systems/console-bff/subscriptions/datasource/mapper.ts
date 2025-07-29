/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { logger } from "../../common/logger";
import { formatKPIValue, isMetricNetworkCheckFailed } from "../../common/utils";
import {
  GetLatestMetricInput,
  GetMetricsStatInput,
  LatestMetricRes,
  MetricRes,
  MetricsRes,
  NotificationsAPIRes,
  NotificationsAPIResDto,
  NotificationsRes,
  NotificationsResDto,
} from "../resolvers/types";
import { eventKeyToAction } from "./../../common/notification/index";

const ERROR_RESPONSE = {
  success: true,
  msg: "success",
  orgId: "",
  nodeId: "",
  type: "",
};

export const parseLatestMetricRes = (
  res: any,
  args: GetLatestMetricInput
): LatestMetricRes => {
  const data = res.data.result[0];
  if (data?.value?.length > 0) {
    return {
      success: true,
      msg: "success",
      nodeId: args.nodeId,
      siteId: data.siteid,
      type: args.type,
      value: data.value,
    };
  } else {
    return { ...ERROR_RESPONSE, value: [0, 0] } as LatestMetricRes;
  }
};

export const parseMetricsResponse = (
  res: Array<{
    metric: {
      env: string;
      nodeid?: string;
      site?: string;
      network?: string;
      package?: string;
      dataplan?: string;
    };
    values: [number, string][];
  }>,
  type: string,
  args: GetMetricsStatInput
): MetricsRes => {
  const metricResArray: MetricRes[] = res.map(item => {
    if (
      isMetricNetworkCheckFailed(
        args.networkId || "",
        item.metric.network || "",
        args.operation || ""
      )
    ) {
      return { ...ERROR_RESPONSE, values: [[0, 0]] };
    }
    if (
      isMetricNetworkCheckFailed(
        args.siteId || "",
        item.metric.site || "",
        args.operation || ""
      )
    ) {
      return { ...ERROR_RESPONSE, values: [[0, 0]] };
    }
    return {
      type: type,
      success: true,
      msg: "success",
      nodeId: item.metric.nodeid ?? args.nodeId ?? "",
      siteId: item.metric?.site ?? args.siteId ?? "",
      networkId: item.metric?.network ?? args.networkId ?? "",
      packageId: item.metric?.package ?? "",
      dataPlanId: item.metric?.dataplan ?? "",
      values: fixTimestampInMetricData(
        item.values,
        1,
        args.to ?? Math.floor(Date.now() / 1000),
        args.from,
        type
      ),
    };
  });

  const metricsRes: MetricsRes = {
    metrics: metricResArray,
  };

  return metricsRes;
};

function fixTimestampInMetricData(
  data: [number, string | null][],
  step: number,
  to: number,
  from: number,
  type: string
): [number, number][] {
  logger.info("fixTimestampInMetricData", data);
  if (!Array.isArray(data) || data.length === 0) return [];
  const result: [number, number][] = [];
  let prevTimestamp: number = from;
  let dataIndex = 0;

  while (prevTimestamp <= to) {
    if (dataIndex < data.length && data[dataIndex][0] === prevTimestamp) {
      result.push([
        data[dataIndex][0] * 1000,
        formatKPIValue(type, data[dataIndex][1]),
      ]);
      dataIndex++;
    } else {
      result.push([prevTimestamp * 1000, -1]);
    }
    prevTimestamp += step;
  }

  return result;
}

export const parseNotification = (
  notification: NotificationsAPIResDto
): NotificationsResDto => {
  const n: NotificationsResDto = {
    id: notification.id,
    type: notification.type,
    scope: notification.scope,
    title: notification.title,
    isRead: notification.is_read,
    eventKey: notification.event_key,
    createdAt: notification.created_at,
    resourceId: notification.resource_id,
    description: notification.description,
  };
  n.redirect = eventKeyToAction(notification.event_key, n);
  return n;
};

export const parseNotificationsRes = (
  res: NotificationsAPIRes
): NotificationsRes => {
  const data = res.notifications;
  const notifications: NotificationsResDto[] = [];
  data.map((notification: NotificationsAPIResDto) => {
    notifications.push(parseNotification(notification));
  });
  return {
    notifications: notifications,
  };
};
