import { Circle, MoreHoriz } from '@mui/icons-material';
import {
  Box,
  Divider,
  IconButton,
  List,
  ListItem,
  Popover,
  Typography,
} from '@mui/material';
import { useState } from 'react';
import DeleteNotification from '../DeleteNotification';

interface Alert {
  id: number;
  isRead: boolean;
  title: string;
  description: string;
  time: string;
}

interface AlertBoxProps {
  alerts: Alert[];
  onAlertRead: (index: number) => void;
}

const AlertBox = ({ alerts, onAlertRead }: AlertBoxProps) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  // const [readAlerts, setReadAlerts] = useState<Set<String>>(new Set());

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  // const handleAlertClick = (title: string) => {
  //   setReadAlerts((prev) => new Set(prev.add(title)));
  // };

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
        <Typography variant="h6" fontWeight="500" fontFamily="Rubik">
          Alerts
        </Typography>
        <Typography
          fontSize="16px"
          fontWeight="lighter"
          fontFamily="Work Sans"
          paddingLeft={1}
        >
          ({alerts.filter((alert) => !alert.isRead).length})
        </Typography>
      </Box>
      <Divider sx={{ margin: 0 }} />
      <List sx={{ padding: 0, margin: 0 }}>
        {alerts.map((alert: Alert, index: number) => (
          <Box key={alert.id} sx={{ margin: 0 }}>
            <ListItem
              alignItems="flex-start"
              sx={{
                bgcolor: alert.isRead ? 'none' : '#007DFF12',
                cursor: 'pointer',
                flexDirection: 'column',
                alignItems: 'flex-start',
              }}
              onClick={() => onAlertRead(index)}
            >
              <Box display="flex" alignItems="center" width="100%">
                {!alert.isRead && (
                  <Circle
                    sx={{ fontSize: '12px', marginRight: 1 }}
                    color="secondary"
                  />
                )}
                <Typography
                  fontSize="16px"
                  fontWeight="500"
                  fontFamily="Work Sans"
                >
                  {alert.title}
                </Typography>
                <Box flexGrow={1} />
                <Typography
                  fontSize="12px"
                  fontWeight="400"
                  fontFamily="Work Sans"
                >
                  {alert.time}
                </Typography>
              </Box>
              <Box display="flex" alignItems="center" width="100%">
                <Typography
                  variant="body2"
                  sx={{ flexGrow: 1 }}
                  fontSize="14px"
                  fontWeight="400"
                  fontFamily="Work Sans"
                >
                  {alert.description}
                </Typography>
                <IconButton onClick={handleMenuClick}>
                  <MoreHoriz />
                </IconButton>
              </Box>
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
                <DeleteNotification />
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
