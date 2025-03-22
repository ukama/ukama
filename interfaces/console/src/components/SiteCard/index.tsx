/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import RouterIcon from '@mui/icons-material/Router';
import BatteryChargingFullIcon from '@mui/icons-material/BatteryChargingFull';
import BatteryAlertIcon from '@mui/icons-material/BatteryAlert';
import Battery50Icon from '@mui/icons-material/Battery50';
import SignalCellularAltIcon from '@mui/icons-material/SignalCellularAlt';
import SignalCellular1BarIcon from '@mui/icons-material/SignalCellular1Bar';
import SignalCellular2BarIcon from '@mui/icons-material/SignalCellular2Bar';
import {
  Box,
  Card,
  CardContent,
  IconButton,
  Menu,
  MenuItem,
  Skeleton,
  Typography,
  useMediaQuery,
} from '@mui/material';
import { useRouter } from 'next/navigation';
import React, { useState } from 'react';

interface SiteCardProps {
  siteId: string;
  name: string;
  address: string;
  userCount?: number;
  connectionStatus?: string;
  batteryStatus?: string;
  signalStrength?: string;
  onClickMenu?: (siteId: string) => void;
  loading?: boolean;
  handleSiteNameUpdate: (siteId: string, newSiteName: string) => void;
}

const SiteCard: React.FC<SiteCardProps> = ({
  siteId,
  name,
  address,
  connectionStatus = 'Online',
  batteryStatus = 'Charged',
  signalStrength = 'Strong',
  onClickMenu,
  handleSiteNameUpdate,
  loading = false,
}) => {
  const router = useRouter();

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
  };
  const handleMenuClick = (action: string) => {
    if (handleSiteNameUpdate) {
      handleSiteNameUpdate(siteId, name);
    }
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  const navigateToDetails = () => {
    router.push(`/console/sites/${siteId}`);
  };

  const getConnectionStyles = () => {
    switch (connectionStatus) {
      case 'Online':
        return {
          color: colors.green,
          icon: <RouterIcon sx={{ color: colors.green }} />,
        };
      case 'Offline':
        return {
          color: colors.red,
          icon: <RouterIcon sx={{ color: colors.red }} />,
        };
      case 'Warning':
        return {
          color: colors.orange,
          icon: <RouterIcon sx={{ color: colors.orange }} />,
        };
      default:
        return {
          color: colors.green,
          icon: <RouterIcon sx={{ color: colors.green }} />,
        };
    }
  };

  const getBatteryStyles = () => {
    switch (batteryStatus) {
      case 'Charged':
        return {
          color: colors.green,
          icon: <BatteryChargingFullIcon sx={{ color: colors.green }} />,
        };
      case 'Medium':
        return {
          color: colors.orange,
          icon: <Battery50Icon sx={{ color: colors.orange }} />,
        };
      case 'Low':
        return {
          color: colors.red,
          icon: <BatteryAlertIcon sx={{ color: colors.red }} />,
        };
      default:
        return {
          color: colors.green,
          icon: <BatteryChargingFullIcon sx={{ color: colors.green }} />,
        };
    }
  };

  const getSignalStyles = () => {
    switch (signalStrength) {
      case 'Strong':
        return {
          color: colors.green,
          icon: <SignalCellularAltIcon sx={{ color: colors.green }} />,
        };
      case 'Medium':
        return {
          color: colors.orange,
          icon: <SignalCellular2BarIcon sx={{ color: colors.orange }} />,
        };
      case 'Weak':
        return {
          color: colors.red,
          icon: <SignalCellular1BarIcon sx={{ color: colors.red }} />,
        };
      default:
        return {
          color: colors.green,
          icon: <SignalCellularAltIcon sx={{ color: colors.green }} />,
        };
    }
  };

  const connectionStyles = getConnectionStyles();
  const batteryStyles = getBatteryStyles();
  const signalStyles = getSignalStyles();

  return (
    <Card
      sx={{
        border: `1px solid ${colors.darkGradient}`,
        borderRadius: 2,
        marginBottom: 2,
        backgroundColor: colors.lightGray,
      }}
      onClick={navigateToDetails}
    >
      <CardContent>
        <Box
          display="flex"
          justifyContent="space-between"
          alignItems="flex-start"
        >
          <Box>
            <Typography
              variant="h6"
              sx={{
                borderBottom: '1px solid',
                display: 'inline-block',
                mb: 1,
                fontWeight: 'bold',
                cursor: 'pointer',
              }}
            >
              {loading ? <Skeleton width={150} /> : name}
            </Typography>

            <Typography color="textSecondary" variant="body1">
              {loading ? <Skeleton width={200} /> : address}
            </Typography>
          </Box>

          <IconButton onClick={handleClick} sx={{ mt: -1 }}>
            <MoreVertIcon />
          </IconButton>

          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleClose}
            onClick={(e) => e.stopPropagation()}
          >
            <MenuItem onClick={() => handleMenuClick('edit')}>
              Edit Name
            </MenuItem>
          </Menu>
        </Box>

        <Box
          display="flex"
          justifyContent="flex-start"
          alignItems="center"
          mt={3}
          gap={4}
        >
          <Box display="flex" alignItems="center" gap={1}>
            {loading ? (
              <Skeleton width={24} height={24} />
            ) : (
              connectionStyles.icon
            )}
            <Typography variant="body2" sx={{ color: connectionStyles.color }}>
              {loading ? <Skeleton width={60} /> : connectionStatus}
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" gap={1}>
            {loading ? <Skeleton width={24} height={24} /> : batteryStyles.icon}
            <Typography variant="body2" sx={{ color: batteryStyles.color }}>
              {loading ? <Skeleton width={70} /> : batteryStatus}
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" gap={1}>
            {loading ? <Skeleton width={24} height={24} /> : signalStyles.icon}
            <Typography variant="body2" sx={{ color: signalStyles.color }}>
              {loading ? <Skeleton width={60} /> : signalStrength}
            </Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default SiteCard;
