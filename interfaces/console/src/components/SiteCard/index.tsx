/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import {
  getBatteryStyles,
  getConnectionStyles,
  getSignalStyles,
} from '@/utils';
import MoreVertIcon from '@mui/icons-material/MoreVert';

import {
  Box,
  Card,
  CardContent,
  IconButton,
  Menu,
  MenuItem,
  Skeleton,
  Typography,
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

  const connectionStyles = getConnectionStyles(connectionStatus);
  const batteryStyles = getBatteryStyles(batteryStatus);
  const signalStyles = getSignalStyles(signalStrength);

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
