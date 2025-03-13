/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import PersonIcon from '@mui/icons-material/Person';
import RouterIcon from '@mui/icons-material/Router';
import BatteryChargingFullIcon from '@mui/icons-material/BatteryChargingFull';
import SignalCellularAltIcon from '@mui/icons-material/SignalCellularAlt';
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
  userCount = 2,
  connectionStatus = 'Online',
  batteryStatus = 'Charged',
  signalStrength = 'Strong',
  onClickMenu,
  handleSiteNameUpdate,
  loading = false,
}) => {
  const isSmallScreen = useMediaQuery((theme: any) =>
    theme.breakpoints.down('sm'),
  );
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

  return (
    <Card
      sx={{
        border: `1px solid ${colors.darkGradient}`,
        borderRadius: 2,
        marginBottom: 2,
        '&:hover': {
          border: `1px solid ${colors.primaryDark}`,
          cursor: 'pointer',
        },
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
                fontWeight: 'bold'
              }}
            >
              {loading ? (
                <Skeleton width={150} />
              ) : (
                name
              )}
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
        
        {/* Status indicators row */}
        <Box
          display="flex"
          justifyContent="flex-start"
          alignItems="center"
          mt={3}
          gap={4}
        >
          {/* User count */}
          <Box display="flex" alignItems="center" gap={1}>
            <PersonIcon color="action" />
            <Typography variant="body2">
              {loading ? <Skeleton width={20} /> : userCount}
            </Typography>
          </Box>
          
          {/* Connection status */}
          <Box display="flex" alignItems="center" gap={1}>
            <RouterIcon sx={{ color: connectionStatus === 'Online' ? 'green' : 'gray' }} />
            <Typography variant="body2" color={connectionStatus === 'Online' ? 'green' : 'textSecondary'}>
              {loading ? <Skeleton width={60} /> : connectionStatus}
            </Typography>
          </Box>
          
          {/* Battery status */}
          <Box display="flex" alignItems="center" gap={1}>
            <BatteryChargingFullIcon sx={{ color: '#d32f2f' }} />
            <Typography variant="body2" color="#d32f2f">
              {loading ? <Skeleton width={70} /> : batteryStatus}
            </Typography>
          </Box>
          
          {/* Signal strength */}
          <Box display="flex" alignItems="center" gap={1}>
            <SignalCellularAltIcon sx={{ color: 'green' }} />
            <Typography variant="body2" color="green">
              {loading ? <Skeleton width={60} /> : signalStrength}
            </Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default SiteCard;