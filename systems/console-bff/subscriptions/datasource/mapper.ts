/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  GetLatestMetricInput,
  GetMetricRangeInput,
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
      orgId: data.metric.org,
      nodeId: args.nodeId,
      type: args.type,
      value: data.value,
    };
  } else {
    return { ...ERROR_RESPONSE, value: [0, 0] } as LatestMetricRes;
  }
};

export const parseMetricRes = (res: any, type: string): MetricRes => {
  const data = res.data.result[0];
  if (data?.values?.length > 0) {
    return {
      type: type,
      success: true,
      msg: "success",
      orgId: data.metric.org,
      nodeId: data.metric.nodeid,
      values: fixTimestampInMetricData(data.values),
    };
  } else {
    return { ...ERROR_RESPONSE, values: [[0, 0]] } as MetricRes;
  }
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
        orgId: result[0].metric.org,
        nodeId: result[0].metric.nodeid,
        values: fixTimestampInMetricData(result[0].values),
      }
    : getEmptyMetric(args);
};

const fixTimestampInMetricData = (
  values: [[number, string]]
): [number, number][] => {
  if (values.length > 0) {
    const fixedValues: [number, number][] = values.map(
      (value: [number, string]) => {
        return [Math.floor(value[0]) * 1000, parseFloat(value[1])];
      }
    );
    return fixedValues;
  }
  return [];
};

export const parsePromethRes = (
  res: any,
  args: GetMetricRangeInput
): MetricRes => {
  const metric = res.data.result.filter(
    (item: any) => item.metric.nodeid === args.nodeId
  )[0];

  if (metric?.values?.length > 0) {
    return {
      type: args.type,
      success: true,
      msg: "success",
      orgId: metric.metric.org,
      nodeId: metric.metric.nodeid,
      values: fixTimestampInMetricData(metric.values),
    };
  } else return getEmptyMetric(args);
};

export const parseNotification = (
  notification: NotificationsAPIResDto
): NotificationsResDto => {
  return {
    id: notification.id,
    type: notification.type,
    scope: notification.scope,
    title: notification.title,
    isRead: notification.is_read,
    createdAt: notification.created_at,
    description: notification.description,
    nodeStateId: notification.node_state_id || "",
    nodeState: notification.node_state
      ? {
          Id: notification.node_state.Id,
          name: notification.node_state.name,
          nodeId: notification.node_state.nodeId,
          current_state: notification.node_state.current_state,
          latitude: notification.node_state.latitude,
          longitude: notification.node_state.longitude,
          created_at: notification.node_state.created_at,
        }
      : null,
  };
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
