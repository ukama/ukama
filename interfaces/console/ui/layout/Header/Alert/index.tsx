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
  const [alerts, setAlerts] = useState([
    {
      id:1,
      title: 'Alert 1',
      description: 'Item affected + severity + type of alert',
      time: '8/30 1PM',
      isRead: false,
    },
    {

      id:2,
      title: 'Alert 2',
      description: 'Item affected + severity + type of alert',
      time: '8/16 1PM',
      isRead: false,
    },
    {
      id:3,
      title: 'Alert 3',
      description: 'Item affected + severity + type of alert',
      time: '8/20 1PM',
      isRead: false,
    },
  ]);

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  const handleAlertRead=(index: number)=>{
    setAlerts((prev)=>{
      const newAlerts = [...prev]
      newAlerts[index] = {...newAlerts[index], isRead:true}
      return newAlerts
    })
  }

  const unreadCount = alerts.filter(alert=>!alert.isRead).length
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
        <AlertBox alerts={alerts} onAlertRead={handleAlertRead}/>
      </Popover>
    </>
  );
};

export default Alerts;
