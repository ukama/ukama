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
    resourceId: res.notification.resource_id,
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
