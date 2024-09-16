/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  NotificationAPIRes,
  NotificationResDto,
  NotificationsAPIRes,
  NotificationsResDto,
  UpdateNotificationResDto,
} from "../resolvers/types";

export const dtoToNotificationDto = (
  res: NotificationAPIRes
): NotificationResDto => {
  return {
    id: res.notification.id,
    title: res.notification.title,
    description: res.notification.description,
    type: res.notification.type,
    scope: res.notification.scope,
    createdAt: res.notification.created_at,
    orgId: res.notification.org_id,
    networkId: res.notification.network_id,
    userId: res.notification.user_id,
    subscriberId: res.notification.subscriber_id,
    nodeStateId: res.notification.node_state_id || "",
    nodeState: res.notification.node_state
      ? {
          Id: res.notification.node_state.Id,
          name: res.notification.node_state.name,
          nodeId: res.notification.node_state.nodeId,
          current_state: res.notification.node_state.current_state,
          latitude: res.notification.node_state.latitude,
          longitude: res.notification.node_state.longitude,
          created_at: res.notification.node_state.created_at,
        }
      : null,
  };
};

export const dtoToNotificationsDto = (
  res: NotificationsAPIRes
): NotificationsResDto => {
  return {
    notifications: res.notifications.map(notification => ({
      id: notification.id,
      title: notification.title,
      description: notification.description,
      type: notification.type,
      scope: notification.scope,
      isRead: notification.is_read,
      createdAt: notification.created_at,
      nodeStateId: notification.node_state_id,
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
    })),
  };
};

export const dtoToUpdateNotificationDto = (
  res: UpdateNotificationResDto
): UpdateNotificationResDto => {
  return {
    id: res.id,
  };
};
