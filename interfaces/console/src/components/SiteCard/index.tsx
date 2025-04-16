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
} from '@mui/material';
import React, { useEffect, useState, useRef, memo } from 'react';
import PubSub from 'pubsub-js';

interface SiteCardProps {
  siteId: string;
  name: string;
  address: string;
  userCount?: number;
  loading?: boolean;
  handleSiteNameUpdate: (siteId: string, newSiteName: string) => void;
  maxAddressLength?: number;
}

const truncateText = (text: string, maxLength: number): string => {
  if (!text || typeof text !== 'string') return '';
  if (text.length <= maxLength) return text;
  return `${text.substring(0, maxLength)}...`;
};

const SiteCard: React.FC<SiteCardProps> = memo(
  ({
    siteId,
    name,
    address,
    userCount = 0,
    handleSiteNameUpdate,
    loading,
    maxAddressLength = 49,
  }) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [metrics, setMetrics] = useState({
      site_uptime_seconds: null,
      battery_charge_percentage: null,
      backhaul_speed: null,
    });

    const subscriptionToken = useRef<string | null>(null);

    const displayAddress = loading
      ? ''
      : truncateText(address || '', maxAddressLength);

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

    useEffect(() => {
      if (!siteId || loading) return;

      const token = PubSub.subscribe(
        `site-metrics-${siteId}`,
        (msg, { type, value }) => {
          setMetrics((prev) => ({
            ...prev,
            [type]: value,
          }));
        },
      );

      subscriptionToken.current = token;

      PubSub.publish(`request-metrics-${siteId}`, {});

      return () => {
        if (subscriptionToken.current) {
          PubSub.unsubscribe(subscriptionToken.current);
          subscriptionToken.current = null;
        }
      };
    }, [siteId, loading]);

    const connectionStyles =
      metrics.site_uptime_seconds !== null
        ? getStatusStyles('uptime', metrics.site_uptime_seconds)
        : { icon: null, color: colors.darkGray };

    const batteryStyles =
      metrics.battery_charge_percentage !== null
        ? getStatusStyles('battery', metrics.battery_charge_percentage)
        : { icon: null, color: colors.darkGray };

    const signalStyles =
      metrics.backhaul_speed !== null
        ? getStatusStyles('signal', metrics.backhaul_speed)
        : { icon: null, color: colors.darkGray };

    const isMetricAvailable = (metric: any) => {
      return !loading && metric !== null;
    };

    return (
      <Card
        sx={{
          border: `1px solid ${colors.darkGradient}`,
          borderRadius: 2,
          marginBottom: 2,
          backgroundColor: colors.lightGray,
          cursor: 'pointer',
          transition: 'transform 0.2s, box-shadow 0.2s',
          '&:hover': {
            transform: 'translateY(-2px)',
            boxShadow: '0 4px 8px rgba(0,0,0,0.1)',
          },
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
                <Tooltip title={address || ''} placement="top-start">
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

            <IconButton onClick={handleClick} sx={{ mt: -1 }} aria-label="menu">
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
                <Skeleton width={24} height={24} variant="rectangular" />
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
              {!isMetricAvailable(metrics.site_uptime_seconds) ? (
                <Skeleton width={24} height={24} variant="circular" />
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
                {!isMetricAvailable(metrics.site_uptime_seconds) ? (
                  <Skeleton width={60} />
                ) : metrics.site_uptime_seconds !== null &&
                  metrics.site_uptime_seconds <= 0 ? (
                  'Offline'
                ) : (
                  'Online'
                )}
              </Typography>
            </Box>

            <Box display="flex" alignItems="center" gap={1}>
              {!isMetricAvailable(metrics.battery_charge_percentage) ? (
                <Skeleton width={24} height={24} variant="circular" />
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
                {!isMetricAvailable(metrics.battery_charge_percentage) ? (
                  <Skeleton width={70} />
                ) : metrics.battery_charge_percentage !== null &&
                  metrics.battery_charge_percentage < 20 ? (
                  'Critical'
                ) : metrics.battery_charge_percentage !== null &&
                  metrics.battery_charge_percentage < 40 ? (
                  'Low'
                ) : (
                  'Charged'
                )}
              </Typography>
            </Box>

            <Box display="flex" alignItems="center" gap={1}>
              {!isMetricAvailable(metrics.backhaul_speed) ? (
                <Skeleton width={24} height={24} variant="circular" />
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
                {!isMetricAvailable(metrics.backhaul_speed) ? (
                  <Skeleton width={60} />
                ) : metrics.backhaul_speed !== null &&
                  metrics.backhaul_speed < 10 ? (
                  'No signal'
                ) : metrics.backhaul_speed !== null &&
                  metrics.backhaul_speed < 70 ? (
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
  },
);

SiteCard.displayName = 'SiteCard';

export default SiteCard;
