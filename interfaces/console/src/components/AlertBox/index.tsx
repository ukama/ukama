import React, { useState } from 'react';
import { NotificationsResDto, NodeStateResDto } from '@/client/graphql/generated/subscriptions';
import { Circle, MoreHoriz } from '@mui/icons-material';
import {
  Box,
  Divider,
  IconButton,
  List,
  ListItem,
  Popover,
  Typography,
  Button,
} from '@mui/material';
import { format } from 'date-fns';

interface AlertBoxProps {
  alerts: NotificationsResDto[] | undefined;
  handleNotificationRead: (id: string) => void;
  configShowButtonState: string;
  onConfigureSite: (nodeState: NodeStateResDto) => void;
}

const AlertBox: React.FC<AlertBoxProps> = ({
  alerts,
  handleNotificationRead,
  configShowButtonState,
  onConfigureSite,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'alert-popover' : undefined;

  return (
    <Box
      bgcolor={'white'}
      borderRadius={'10px'}
      width={'398px'}
      height={'310px'}
    >
      <Box display="flex" justifyContent="flex-start" alignItems="center" p={2}>
        <Typography variant="h6">Alerts</Typography>
        <Typography variant="body1" paddingLeft={1}>
          ({alerts?.filter((alert) => !alert.isRead).length ?? 0})
        </Typography>
      </Box>
      <Divider sx={{ margin: 0 }} />
      <List sx={{ padding: 0, margin: 0 }}>
        {alerts?.map((alert: NotificationsResDto) => (
          <Box key={alert.id} sx={{ margin: 0 }}>
            <ListItem
              alignItems="flex-start"
              sx={{
                bgcolor: alert.isRead ? 'none' : '#007DFF12',
                cursor: 'pointer',
                flexDirection: 'column',
                alignItems: 'flex-start',
              }}
              onClick={() => handleNotificationRead(alert.id)}
            >
              <Box display="flex" alignItems="center" width="100%">
                {!alert.isRead && (
                  <Circle
                    sx={{ fontSize: '12px', marginRight: 1 }}
                    color="secondary"
                  />
                )}
                <Typography fontSize="16px" fontWeight="500">
                  {alert.title}
                </Typography>
                <Box flexGrow={1} />
                <Typography fontSize="12px" fontWeight="400">
                  {format(new Date(alert.createdAt), 'MM/dd hh:mm a')}
                </Typography>
              </Box>
              <Box display="flex" alignItems="center" width="100%">
                <Typography variant="body2" sx={{ flexGrow: 1 }}>
                  {alert.description}
                </Typography>
                <IconButton onClick={handleMenuClick}>
                  <MoreHoriz />
                </IconButton>
              </Box>
              {alert.nodeState &&
                alert.nodeState.currentState === configShowButtonState && (
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={() => {
                      if (alert.nodeState) {
                        onConfigureSite(alert.nodeState);
                      }
                    }}
                    sx={{ mt: 1 }}
                  >
                    Configure Site
                  </Button>
                )}
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
                {/* <DeleteNotification /> */}
              </Popover>
            </ListItem>
            <Divider sx={{ margin: 0 }} />
          </Box>
        ))}
      </List>
    </Box>
  );
};

export default AlertBox;