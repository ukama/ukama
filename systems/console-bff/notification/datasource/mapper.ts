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
          id: res.notification.node_state.id,
          name: res.notification.node_state.name,
          node_id: res.notification.node_state.node_id,
          currentState: res.notification.node_state.currentState,
          latitude: res.notification.node_state.latitude,
          longitude: res.notification.node_state.longitude,
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
      is_read: notification.is_read,
      created_at: notification.created_at,
      nodeStateId: notification.nodeStateId,
      nodeState: notification.nodeState
        ? {
            id: notification.nodeState.id,
            name: notification.nodeState.name,
            node_id: notification.nodeState.node_id,
            latitude: notification.nodeState.latitude,
            longitude: notification.nodeState.longitude,
            currentState: notification.nodeState.currentState,
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
