import {
  NotificationRes,
  useGetNotificationsQuery,
  useGetNotificationsSubSubscription,
} from '@/generated/metrics';
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

const Alerts = () => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [alerts, setAlerts] = useState<NotificationRes[] | undefined>(
    undefined,
  );

  // Fetch initial notifications
  const { data: queryData } = useGetNotificationsQuery({
    fetchPolicy: 'cache-and-network',
    variables: {
      data: {
        orgId: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b221',
        userId: 'da421ed5-0fba-4638-9661-a9204f49006a',
        networkId: 'da421ed5-0fba-4638-9661-a9204f490069',
        scopes: ['notifications'],
        siteId: 'da421ed5-0fba-4638-9661-a9204f490062',
        subscriberId: 'da421ed5-0fba-4638-9661-a9204f490065',
      },
    },
    onCompleted: (data) => {
      console.log('Query completed:', data);
      setAlerts(data.getNotifications.notifications);
    },
  });

  // Subscribe to notifications
  console.log(useGetNotificationsSubSubscription({
    variables: {
      orgId: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b221',
      userId: 'da421ed5-0fba-4638-9661-a9204f49006a',
      networkId: 'da421ed5-0fba-4638-9661-a9204f490069',
      scopes: ['notifications'],
      siteId: 'da421ed5-0fba-4638-9661-a9204f490062',
      subscriberId: 'da421ed5-0fba-4638-9661-a9204f490065',
    },
    onData: ({ data: subscriptionData }) => {
      const newAlerts = subscriptionData.data?.getNotificationsSub;
      console.log('Subscription data:', newAlerts);
      if (newAlerts) {
        setAlerts((prev) => (prev ? [newAlerts, ...prev] : [newAlerts]));
      }
    },
  }))

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
    setAlerts((prev) => {
      if (!prev) return prev;
      const newAlerts = [...prev];
      newAlerts[index] = { ...newAlerts[index], isRead: true };
      return newAlerts;
    });
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
