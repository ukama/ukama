/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NotificationsRes } from '@/client/graphql/generated/subscriptions';
import { Notifications } from '@/components/NotificationsUI';
import { IconStyle } from '@/styles/global';
import NotificationsIcon from '@mui/icons-material/Notifications';
import { Badge, IconButton, Popover } from '@mui/material';
import React, { useState } from 'react';

interface IAlertsProps {
  notifications: NotificationsRes;
  handleNotificationRead: (id: string) => void;
}

const Alerts = ({ notifications, handleNotificationRead }: IAlertsProps) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  // Handle popover open
  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  // Handle popover close
  const handleClose = () => {
    setAnchorEl(null);
  };

  const unreadCount = notifications.notifications?.filter(
    (alert) => !alert.isRead,
  ).length;
  const open = Boolean(anchorEl);
  const id = open ? 'alert-popover' : undefined;

  return (
    <>
      <IconButton sx={{ ...IconStyle }} onClick={handleClick}>
        <Badge badgeContent={unreadCount} color="secondary">
          <NotificationsIcon />
        </Badge>
      </IconButton>
      <Popover
        sx={{ mt: 2 }}
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
      >
        <Notifications
          notifications={notifications}
          handleAction={handleNotificationRead}
        />
      </Popover>
    </>
  );
};

export default Alerts;
