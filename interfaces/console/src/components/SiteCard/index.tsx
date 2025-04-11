/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import { getStatusStyles } from '@/utils';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import PeopleIcon from '@mui/icons-material/People';
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
import React, { useState } from 'react';

interface SiteCardProps {
  siteId: string;
  name: string;
  address: string;
  userCount?: number;
  siteUptimeSeconds?: number | null;
  batteryPercentage?: number | null;
  backhaulSpeed?: number | null;
  loading?: boolean;
  handleSiteNameUpdate: (siteId: string, newSiteName: string) => void;
}

const SiteCard: React.FC<SiteCardProps> = ({
  siteId,
  name,
  address,
  userCount = 0,
  siteUptimeSeconds,
  batteryPercentage,
  backhaulSpeed,
  handleSiteNameUpdate,
  loading = false,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleMenuClick = () => {
    handleSiteNameUpdate(siteId, name);
    handleClose();
  };

  const navigateToDetails = () => {
    window.location.href = `/console/sites/${siteId}`;
  };

  const connectionStyles = getStatusStyles('uptime', siteUptimeSeconds ?? 0);
  const batteryStyles = getStatusStyles('battery', batteryPercentage ?? 0);
  const signalStyles = getStatusStyles('signal', backhaulSpeed ?? 0);

  return (
    <Card
      sx={{
        border: `1px solid ${colors.darkGradient}`,
        borderRadius: 2,
        marginBottom: 2,
        backgroundColor: colors.lightGray,
        cursor: 'pointer',
      }}
      onClick={navigateToDetails}
    >
      <CardContent>
        <Box display="flex" justifyContent="space-between">
          <Box>
            <Typography
              variant="h6"
              sx={{
                borderBottom: '1px solid',
                display: 'inline-block',
                mb: 1,
                fontWeight: 'bold',
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
            <MenuItem onClick={handleMenuClick}>Edit Name</MenuItem>
          </Menu>
        </Box>

        <Box display="flex" mt={3} gap={4}>
          <Box display="flex" alignItems="center" gap={1}>
            {loading ? (
              <Skeleton width={24} height={24} />
            ) : (
              <PeopleIcon sx={{ color: colors.darkGray }} />
            )}
            <Typography variant="body2" sx={{ color: colors.darkGray }}>
              {loading ? <Skeleton width={30} /> : userCount}
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" gap={1}>
            {loading ? (
              <Skeleton width={24} height={24} />
            ) : (
              connectionStyles.icon
            )}
            <Typography variant="body2" sx={{ color: connectionStyles.color }}>
              {loading ? (
                <Skeleton width={60} />
              ) : (siteUptimeSeconds ?? 0) <= 0 ? (
                'Offline'
              ) : (
                'Online'
              )}
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" gap={1}>
            {loading ? <Skeleton width={24} height={24} /> : batteryStyles.icon}
            <Typography variant="body2" sx={{ color: batteryStyles.color }}>
              {loading ? (
                <Skeleton width={70} />
              ) : (batteryPercentage ?? 0) < 20 ? (
                'Critical'
              ) : (batteryPercentage ?? 0) < 60 ? (
                'Low'
              ) : (
                'Charged'
              )}
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" gap={1}>
            {loading ? <Skeleton width={24} height={24} /> : signalStyles.icon}
            <Typography variant="body2" sx={{ color: signalStyles.color }}>
              {loading ? (
                <Skeleton width={60} />
              ) : (backhaulSpeed ?? 0) < 10 ? (
                'No signal'
              ) : (backhaulSpeed ?? 0) < 70 ? (
                'Low signal'
              ) : (
                'Strong'
              )}
            </Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default SiteCard;
