import { Notification_Type } from '@/client/graphql/generated';
import { NotificationsRes } from '@/client/graphql/generated/subscriptions';
import { colors } from '@/theme';
import { MoreHoriz } from '@mui/icons-material';
import ErrorIcon from '@mui/icons-material/ErrorOutline';
import InfoIcon from '@mui/icons-material/Info';
import CriticalIcon from '@mui/icons-material/ReportGmailerrorred';
import WarningIcon from '@mui/icons-material/WarningAmber';
import {
  Box,
  Divider,
  IconButton,
  List,
  ListItem,
  Popover,
  Stack,
  Typography,
} from '@mui/material';
import { format } from 'date-fns';
import { useState } from 'react';

interface INotifications {
  notifications: NotificationsRes;
  handleAction: (action: string, id: string) => void;
}

export const Notifications = ({
  notifications,
  handleAction,
}: INotifications) => {
  return (
    <Stack
      spacing={1}
      width={'398px'}
      height={'306px'}
      borderRadius={'12px'}
      bgcolor={colors.white}
    >
      <Stack
        p={1.5}
        spacing={0.2}
        direction={'row'}
        alignItems="center"
        justifyContent="flex-start"
      >
        <Typography variant="h6">Notifications</Typography>
        <Typography variant="body1" paddingLeft={1}>
          (
          {notifications.notifications.filter((notif) => !notif.isRead)
            .length ?? 0}
          )
        </Typography>
      </Stack>
      <Divider
        sx={{
          height: '1px',
          width: '100%',
          margin: '0px !important',
          bgcolor: colors.whiteLilac,
        }}
      />
      <List disablePadding sx={{ margin: '0px !important' }}>
        {notifications.notifications.map((notif) => (
          <Stack key={notif.id} spacing={1}>
            <ListItem disablePadding>
              <Notification
                id={notif.id}
                type={notif.type}
                title={notif.title}
                createdAt={notif.createdAt}
                description={notif.description}
                handleMenuItemClick={handleAction}
              />
            </ListItem>
            <Divider
              variant="fullWidth"
              sx={{
                margin: '0px !important',
                bgcolor: colors.whiteLilac,
              }}
            />
          </Stack>
        ))}
      </List>
    </Stack>
  );
};

interface INotification {
  id: string;
  title: string;
  createdAt: string;
  description: string;
  type: Notification_Type;
  handleMenuItemClick: (action: string, id: string) => void;
}

export const Notification = ({
  id: nid,
  type,
  title,
  createdAt,
  description,
  handleMenuItemClick,
}: INotification) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'notification-popover' : undefined;

  return (
    <Stack
      p={1.5}
      spacing={1}
      width={'100%'}
      bgcolor={hexToRGB(GetColorByType(type), 0.2)}
    >
      <Stack direction={'row'} spacing={0.5}>
        {GetIcon(type)}
        <Typography variant="h6" fontSize={'16px'} fontWeight="500">
          {title}
        </Typography>
        <Box flexGrow={1} />
        <Typography fontSize="12px" fontWeight="400">
          {format(new Date(createdAt), 'MM/dd hha')}
        </Typography>
      </Stack>
      <Stack direction={'row'} alignItems={'center'}>
        <Typography variant="body2" fontWeight="400" sx={{ flexGrow: 1 }}>
          {description}
        </Typography>
        <IconButton onClick={handleMenuClick} sx={{ p: 0 }}>
          <MoreHoriz />
        </IconButton>
      </Stack>
      {Action(type)}
      <Popover
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
      >
        <Typography
          py={1}
          px={1.5}
          variant="body2"
          onClick={() => handleMenuItemClick('mark-read', nid)}
        >
          Mark as read
        </Typography>
      </Popover>
    </Stack>
  );
};

const Action = (type: Notification_Type) => {
  switch (type) {
    case Notification_Type.NotifActionableInfo:
      return <div></div>;
    case Notification_Type.NotifActionableWarning:
      return <div></div>;
    case Notification_Type.NotifActionableError:
      return <div></div>;
    case Notification_Type.NotifActionableCritical:
      return <div></div>;
    default:
      return <Box display={'none'} />;
  }
};

const GetIcon = (type: Notification_Type) => {
  switch (type) {
    case Notification_Type.NotifInfo:
    case Notification_Type.NotifActionableInfo:
      return <InfoIcon sx={{ svg: { fill: colors.hoverColor } }} />;
    case Notification_Type.NotifWarning:
    case Notification_Type.NotifActionableWarning:
      return <WarningIcon htmlColor={colors.yellow} />;
    case Notification_Type.NotifError:
    case Notification_Type.NotifActionableError:
      return <ErrorIcon htmlColor={colors.redMatt} />;
    case Notification_Type.NotifCritical:
    case Notification_Type.NotifActionableCritical:
      return <CriticalIcon htmlColor={colors.red} />;
    default:
      return <InfoIcon />;
  }
};

const GetColorByType = (type: Notification_Type) => {
  switch (type) {
    case Notification_Type.NotifInfo:
    case Notification_Type.NotifActionableInfo:
      return colors.hoverColor;
    case Notification_Type.NotifWarning:
    case Notification_Type.NotifActionableWarning:
      return colors.yellow;
    case Notification_Type.NotifError:
    case Notification_Type.NotifActionableError:
      return colors.redMatt;
    case Notification_Type.NotifCritical:
    case Notification_Type.NotifActionableCritical:
      return colors.red;
    default:
      return colors.hoverColor;
  }
};

const hexToRGB = (hex: string, alpha: number): string => {
  const h = '0123456789ABCDEF';
  const r = h.indexOf(hex[1]) * 16 + h.indexOf(hex[2]);
  const g = h.indexOf(hex[3]) * 16 + h.indexOf(hex[4]);
  const b = h.indexOf(hex[5]) * 16 + h.indexOf(hex[6]);
  if (alpha) {
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  }

  return `rgba(${r}, ${g}, ${b})`;
};
