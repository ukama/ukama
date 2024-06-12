import { useUpdateNotificationMutation } from '@/generated';
import { NotificationsResDto } from '@/generated/metrics';
import { colors } from '@/styles/theme';
import AlertBox from '@/ui/molecules/AlertBox';
import NotificationsIcon from '@mui/icons-material/Notifications';
import { Badge, IconButton, Popover } from '@mui/material';
import React, { useState } from 'react';

const IconStyle = {
  '.MuiSvgIcon-root': {
    width: '24px',
    height: '24px',
    fill: colors.white,
  },
  '.MuiBadge-root': {
    '.MuiSvgIcon-root': {
      width: '24px',
      height: '24px',
      fill: colors.white,
    },
  },
};

interface IAlertsProps {
  alerts: NotificationsResDto[] | undefined;
  setAlerts: Function;
}

const Alerts = ({ alerts, setAlerts }: IAlertsProps) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [updateNotificationMutation] = useUpdateNotificationMutation();

  // Handle popover open
  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  // Handle popover close
  const handleClose = () => {
    setAnchorEl(null);
  };

  // Mark alert as read
  const handleAlertRead = (index: number) => {
    if (alerts) {
      let alertId = alerts[index].id;
      updateNotificationMutation({
        variables: {
          updateNotificationId: alertId,
          isRead: true,
        },
      });
      setAlerts((prev: any) => {
        if (!prev) return prev;
        const newAlerts = [...prev];
        newAlerts[index] = { ...newAlerts[index], isRead: true };
        return newAlerts;
      });
    }
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
