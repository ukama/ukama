/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { METRICS_INTERVAL } from "../../common/configs";
import { logger } from "../../common/logger";
import { formatKPIValue } from "../../common/utils";
import {
  GetLatestMetricInput,
  GetMetricRangeInput,
  GetMetricsStatInput,
  LatestMetricRes,
  MetricRes,
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

const getEmptyMetric = (args: GetMetricRangeInput): MetricRes => {
  return {
    ...ERROR_RESPONSE,
    type: args.type,
    nodeId: args.nodeId,
    values: [[0, 0]],
  } as MetricRes;
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

export const parseMetricRes = (
  res: any,
  type: string,
  args: GetMetricsStatInput | GetMetricRangeInput
): MetricRes => {
  const { result } = res.data.data;
  logger.info(`Parsing metric response for type ${type}:`, result);

  const hasValues = result.length > 0 && result[0]?.values?.length > 0;
  return hasValues
    ? {
        type: type,
        success: true,
        msg: "success",
        nodeId: result[0].metric?.nodeid ?? "",
        siteId: result[0].metric?.siteid ?? "",
        values: fixTimestampInMetricData(
          result[0].values,
          METRICS_INTERVAL,
          args.to || Date.now(),
          args.from,
          type
        ),
      }
    : getEmptyMetric({
        nodeId: result[0].metric.nodeid,
        orgId: "",
        to: args.to,
        type: args.type,
        siteId: result[0].metric.siteid,
        from: args.from,
        step: args.step,
        userId: args.userId,
        withSubscription: false,
      });
};
export const parseSiteMetricRes = (
  res: any,
  type: string,
  args: GetMetricRangeInput
): MetricRes => {
  let result: any[] = [];
  if (res.data?.data?.result) {
    result = res.data.data.result;
  } else if (res.data?.result) {
    result = res.data.result;
  } else if (res.data) {
    result = Array.isArray(res.data) ? res.data : [res.data];
  } else {
    result = res.result || (Array.isArray(res) ? res : [res]);
  }

  if (!Array.isArray(result)) {
    logger.error(`Unexpected result structure:`, result);
    result = [];
  }

  const hasValues = result.length > 0 && result[0]?.values?.length > 0;

  let siteId = "";
  if (result.length > 0 && result[0]?.metric) {
    siteId = result[0].metric.site || result[0].metric.siteid || "";
  }

  return hasValues
    ? {
        type: type,
        success: true,
        msg: "success",
        siteId: siteId,
        values: fixTimestampInMetricData(
          result[0].values,
          METRICS_INTERVAL,
          args.to || Date.now(),
          args.from,
          type
        ),
      }
    : getEmptyMetric({
        orgId: "",
        to: args.to,
        type: args.type,
        siteId: siteId || args.siteId || "",
        from: args.from,
        step: args.step,
        userId: args.userId,
        withSubscription: false,
      });
};
function fixTimestampInMetricData(
  data: [number, string | null][],
  step: number,
  to: number,
  from: number,
  type: string
): [number, number][] {
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
      result.push([prevTimestamp * 1000, 0]);
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
