/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NotificationsResDto } from '@/client/graphql/generated/metrics';
import AlertBox from '@/components/AlertBox';
import { IconStyle } from '@/styles/global';
import NotificationsIcon from '@mui/icons-material/Notifications';
import { Badge, IconButton, Popover } from '@mui/material';
import React, { useState } from 'react';

interface IAlertsProps {
  alerts: NotificationsResDto[] | undefined;
  handleAlertRead: (index: number) => void;
}

const Alerts = ({ alerts, handleAlertRead }: IAlertsProps) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  // Handle popover open
  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  // Handle popover close
  const handleClose = () => {
    setAnchorEl(null);
  };

  const unreadCount = alerts?.filter((alert) => !alert.isRead).length;
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
        sx={{ mt: 4 }}
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
        <AlertBox alerts={alerts} onAlertRead={handleAlertRead} />
      </Popover>
    </>
  );
};

export default Alerts;
