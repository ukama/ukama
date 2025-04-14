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
  Tooltip,
  Typography,
  useTheme,
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
  maxAddressLength?: number;
}

const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text;
  return `${text.substring(0, maxLength)}...`;
};

const SiteCard: React.FC<SiteCardProps> = ({
  siteId,
  name,
  address,
  userCount = 0,
  siteUptimeSeconds,
  batteryPercentage,
  backhaulSpeed,
  handleSiteNameUpdate,
  loading,
  maxAddressLength = 49,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const displayAddress = loading ? '' : truncateText(address, maxAddressLength);

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
            {loading ? (
              <Typography color="textSecondary" variant="body1">
                <Skeleton width={200} />
              </Typography>
            ) : (
              <Tooltip title={address} placement="top-start">
                <Typography
                  color="textSecondary"
                  variant="body1"
                  sx={{
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    maxWidth: '100%',
                  }}
                >
                  {displayAddress}
                </Typography>
              </Tooltip>
            )}
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
            {loading || userCount === undefined || userCount === null ? (
              <Skeleton width={24} height={24} />
            ) : (
              <PeopleIcon sx={{ color: colors.darkGray }} />
            )}
            <Typography variant="body2" sx={{ color: colors.darkGray }}>
              {loading || userCount === undefined || userCount === null ? (
                <Skeleton width={30} />
              ) : (
                userCount
              )}
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" gap={1}>
            {loading ||
            siteUptimeSeconds == null ||
            siteUptimeSeconds === undefined ? (
              <Skeleton width={24} height={24} />
            ) : (
              connectionStyles.icon
            )}
            <Typography
              variant="body2"
              sx={{
                color: connectionStyles.color,
                display: { xs: 'none', sm: 'block' },
              }}
            >
              {loading ||
              siteUptimeSeconds == null ||
              siteUptimeSeconds === undefined ? (
                <Skeleton width={60} />
              ) : siteUptimeSeconds <= 0 ? (
                'Offline'
              ) : (
                'Online'
              )}
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" gap={1}>
            {loading ||
            batteryPercentage == null ||
            batteryPercentage === undefined ? (
              <Skeleton width={24} height={24} />
            ) : (
              batteryStyles.icon
            )}
            <Typography
              variant="body2"
              sx={{
                color: batteryStyles.color,
                display: { xs: 'none', sm: 'block' },
              }}
            >
              {loading ||
              batteryPercentage == null ||
              batteryPercentage === undefined ? (
                <Skeleton width={70} />
              ) : batteryPercentage < 20 ? (
                'Critical'
              ) : batteryPercentage < 40 ? (
                'Low'
              ) : (
                'Charged'
              )}
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" gap={1}>
            {loading || backhaulSpeed == null || backhaulSpeed === undefined ? (
              <Skeleton width={24} height={24} />
            ) : (
              signalStyles.icon
            )}
            <Typography
              variant="body2"
              sx={{
                color: signalStyles.color,
                display: { xs: 'none', sm: 'block' },
              }}
            >
              {loading ||
              backhaulSpeed == null ||
              backhaulSpeed === undefined ? (
                <Skeleton width={60} />
              ) : backhaulSpeed < 10 ? (
                'No signal'
              ) : backhaulSpeed < 70 ? (
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
