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
  const alerts = [
    {
      title: "Alert 1",
      message: 'Item affected + severity + type of alert',
      time: '8/30 1PM',
      severity: 'mild',
    },
    {
      title: "Alert 2",
      message: 'Item affected + severity + type of alert',
      time: '8/16 1PM',
      severity: 'critical',
    },
    {
      title: "Alert 3",
      message: 'Item affected + severity + type of alert',
      time: '8/20 1PM',
    },
  ];

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'alert-popover' : undefined;

  return (
    <>
      <IconButton sx={{ ...IconStyle }} onClick={handleClick}>
        <Badge badgeContent={alerts.length} color="secondary">
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
        <AlertBox alerts={alerts} />
      </Popover>
    </>
  );
};

export default Alerts;
