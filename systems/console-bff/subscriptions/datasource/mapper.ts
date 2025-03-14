/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { METRICS_INTERVAL } from "../../common/configs";
import { eventKeyToAction, formatKPIValue } from "../../common/utils";
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
  const hasValues = result.length > 0 && result[0]?.values?.length > 0;
  return hasValues
    ? {
        type: type,
        success: true,
        msg: "success",
        nodeId: result[0].metric.nodeid,
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
        from: args.from,
        step: args.step,
        userId: args.userId,
        withSubscription: false,
      });
};

export const parseNodeMetricRes = (
  { code, data }: { code: number; data: any },
  args: GetMetricRangeInput
): MetricRes => {
  if (code === 404) return getEmptyMetric(args);
  const { result } = data.data;
  const hasValues = result.length > 0 && result[0]?.values?.length > 0;

  return hasValues
    ? {
        type: args.type,
        success: true,
        msg: "success",
        nodeId: result[0].metric.nodeid,
        values: fixTimestampInMetricData(
          result[0].values,
          args.step || METRICS_INTERVAL,
          args.to || Date.now(),
          args.from,
          args.type
        ),
      }
    : getEmptyMetric(args);
};
export const parseSiteMetricRes = (
  { code, data }: { code: number; data: any },
  args: GetMetricRangeInput
): MetricRes => {
  if (code === 404 || !data?.data?.result) {
    return getEmptyMetric(args);
  }

  const results = data.data.result;
  if (!Array.isArray(results) || results.length === 0) {
    return getEmptyMetric(args);
  }

  const metricResult =
    results.find(r => {
      const metricName =
        r.metric?.type || r.metric?.name || r.metric?.__name__ || "";
      return metricName.includes(args.type);
    }) || results[0];

  if (!metricResult?.values) {
    return getEmptyMetric(args);
  }

  try {
    const values = metricResult.values.map(
      (pair: [number | string, string | number]) => {
        const timestamp =
          typeof pair[0] === "string" ? parseInt(pair[0], 10) : pair[0];
        const value =
          typeof pair[1] === "string" ? parseFloat(pair[1]) : pair[1];
        return [timestamp * 1000, value] as [number, number];
      }
    );

    const metrics: MetricRes = {
      type: args.type,
      success: true,
      msg: "success",
      siteId: metricResult.metric?.site || args.siteId || "",
      values: fixTimestampInMetricData(
        values,
        args.step || METRICS_INTERVAL,
        args.to || Date.now(),
        args.from,
        args.type
      ),
    };

    return metrics;
  } catch (error) {
    return getEmptyMetric(args);
  }
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
