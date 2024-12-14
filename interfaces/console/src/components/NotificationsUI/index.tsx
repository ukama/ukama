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
  Button,
  Divider,
  IconButton,
  Link,
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
      minHeight={'264px'}
      maxHeight={'300px'}
      borderRadius={'12px'}
      height={'fit-contnet'}
      bgcolor={colors.white}
    >
      <Stack
        py={1.2}
        px={1.5}
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

      <List disablePadding sx={{ margin: '0px !important', height: '100%' }}>
        {notifications.notifications.length > 0 ? (
          notifications.notifications.map((notif) => (
            <Stack key={notif.id} spacing={1}>
              <ListItem disablePadding>
                <Notification
                  id={notif.id}
                  type={notif.type}
                  title={notif.title}
                  isRead={notif.isRead}
                  createdAt={notif.createdAt}
                  action={notif.redirect.action}
                  description={notif.description}
                  actionTitle={notif.redirect.title}
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
          ))
        ) : (
          <Stack py={1} alignItems={'center'}>
            <Typography variant="body1" height={'100%'}>
              No notification yet
            </Typography>
          </Stack>
        )}
      </List>
    </Stack>
  );
};

interface INotification {
  id: string;
  title: string;
  action: string;
  isRead: boolean;
  createdAt: string;
  actionTitle: string;
  description: string;
  type: Notification_Type;
  handleMenuItemClick: (action: string, id: string) => void;
}

export const Notification = ({
  id: nid,
  type,
  title,
  isRead,
  action,
  createdAt,
  description,
  actionTitle,
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
      py={1.2}
      px={1.5}
      spacing={0.4}
      width={'100%'}
      bgcolor={isRead ? 'none' : hexToRGB(GetColorByType(type), 0.15)}
    >
      <Stack direction={'row'} spacing={0.8} alignItems={'flex-start'}>
        {GetIcon(type)}
        <Stack
          flexGrow={1}
          spacing={0.5}
          direction={'column'}
          alignContent={'flex-start'}
        >
          <Typography variant="h6" fontSize={'16px'} fontWeight="500">
            {title}
          </Typography>

          <Typography variant="body2" fontWeight="400" sx={{ flexGrow: 1 }}>
            {description}
          </Typography>

          <NotificationAction
            type={type}
            title={actionTitle}
            link={action}
            handleNotificationRead={() => handleMenuItemClick('mark-read', nid)}
          />
        </Stack>

        <Stack direction={'column'} spacing={0.5} alignItems={'flex-end'}>
          <Typography fontSize="12px" fontWeight="400">
            {format(new Date(createdAt), 'MM/dd hha')}
          </Typography>

          <IconButton onClick={handleMenuClick} sx={{ p: 0 }}>
            <MoreHoriz />
          </IconButton>
        </Stack>
      </Stack>

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
        <Button
          variant="text"
          sx={{ textTransform: 'none' }}
          onClick={() => {
            handleMenuItemClick('mark-read', nid);
            handleClose();
          }}
        >
          Mark as read
        </Button>
      </Popover>
    </Stack>
  );
};

interface INotificationAction {
  link: string;
  title: string;
  type: Notification_Type;
  handleNotificationRead: () => void;
}

const NotificationAction = ({
  link,
  type,
  title,
  handleNotificationRead,
}: INotificationAction) => {
  if (
    type !== Notification_Type.TypeActionableInfo &&
    type !== Notification_Type.TypeActionableWarning &&
    type !== Notification_Type.TypeActionableError &&
    type !== Notification_Type.TypeActionableCritical
  ) {
    return <Box display={'none'} />;
  }
  return (
    <Link
      href={link}
      sx={{
        fontSize: '14px',
        fontWeight: 500,
        marginTop: '12px !important',
        color: GetColorByType(type),
      }}
      onClick={handleNotificationRead}
    >
      {title}
    </Link>
  );
};

const GetIcon = (type: Notification_Type) => {
  switch (type) {
    case Notification_Type.TypeInfo:
    case Notification_Type.TypeActionableInfo:
      return (
        <InfoIcon
          sx={{
            fontSize: '20px',
            color: colors.hoverColor,
            marginTop: '4px !important',
          }}
        />
      );
    case Notification_Type.TypeWarning:
    case Notification_Type.TypeActionableWarning:
      return (
        <WarningIcon
          sx={{
            fontSize: '20px',
            color: colors.yellow,
            marginTop: '4px !important',
          }}
        />
      );
    case Notification_Type.TypeError:
    case Notification_Type.TypeActionableError:
      return (
        <ErrorIcon
          sx={{
            fontSize: '20px',
            color: colors.redMatt,
            marginTop: '4px !important',
          }}
        />
      );
    case Notification_Type.TypeCritical:
    case Notification_Type.TypeActionableCritical:
      return (
        <CriticalIcon
          sx={{
            color: colors.red,
            fontSize: '20px',
            marginTop: '4px !important',
          }}
        />
      );
    default:
      return (
        <InfoIcon
          sx={{
            fontSize: '20px',
            color: colors.hoverColor,
            marginTop: '4px !important',
          }}
        />
      );
  }
};

const GetColorByType = (type: Notification_Type) => {
  switch (type) {
    case Notification_Type.TypeInfo:
    case Notification_Type.TypeActionableInfo:
      return colors.hoverColor;
    case Notification_Type.TypeWarning:
    case Notification_Type.TypeActionableWarning:
      return colors.yellow;
    case Notification_Type.TypeError:
    case Notification_Type.TypeActionableError:
      return colors.redMatt;
    case Notification_Type.TypeCritical:
    case Notification_Type.TypeActionableCritical:
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
