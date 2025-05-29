/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { SiteMetricsStateRes } from '@/client/graphql/generated/subscriptions';
import { SITE_KPI_TYPES } from '@/constants';
import colors from '@/theme/colors';
import { getStatusStyles } from '@/utils';
import BatteryAlertIcon from '@mui/icons-material/BatteryAlert';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import RouterIcon from '@mui/icons-material/Router';
import SignalCellularOffIcon from '@mui/icons-material/SignalCellularOff';
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
import PubSub from 'pubsub-js';
import React, { memo, useCallback, useEffect, useRef, useState } from 'react';

const extractMetricValue = (value: any): number | null => {
  if (Array.isArray(value) && value.length > 1) {
    return typeof value[1] === 'number' ? value[1] : null;
  }
  return typeof value === 'number' ? value : null;
};

interface SiteCardProps {
  siteId: string;
  name: string;
  address: string;
  userCount?: number;
  loading?: boolean;
  handleSiteNameUpdate: (siteId: string, newSiteName: string) => void;
  maxAddressLength?: number;
  metricsData?: SiteMetricsStateRes;
}

const truncateText = (text: string, maxLength: number): string => {
  if (!text || typeof text !== 'string') return '';
  if (text.length <= maxLength) return text;
  return `${text.substring(0, maxLength)}...`;
};

const getSiteMetricValue = (
  metricId: string,
  metricsData?: SiteMetricsStateRes,
  siteId?: string,
): number | null => {
  if (!metricsData || !metricsData.metrics || !siteId) return null;

  const metric = metricsData.metrics.find(
    (m) => m.type === metricId && m.success === true && m.siteId === siteId,
  );

  return metric ? extractMetricValue(metric.value) : null;
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
    metricsData,
  }) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

    const [uptimeValue, setUptimeValue] = useState<number | null>(null);
    const [batteryValue, setBatteryValue] = useState<number | null>(null);
    const [backhaulValue, setBackhaulValue] = useState<number | null>(null);

    const subscriptionsRef = useRef<Record<string, string>>({});

    useEffect(() => {
      if (
        !metricsData ||
        !metricsData.metrics ||
        metricsData.metrics.length === 0
      )
        return;

      const uptime = getSiteMetricValue(
        SITE_KPI_TYPES.SITE_UPTIME,
        metricsData,
        siteId,
      );
      if (uptime !== null) setUptimeValue(uptime);

      const batteryCharge = getSiteMetricValue(
        SITE_KPI_TYPES.BATTERY_CHARGE_PERCENTAGE,
        metricsData,
        siteId,
      );
      if (batteryCharge !== null) setBatteryValue(batteryCharge);

      const backhaul = getSiteMetricValue(
        SITE_KPI_TYPES.BACKHAUL_SPEED,
        metricsData,
        siteId,
      );
      if (backhaul !== null) setBackhaulValue(backhaul);
    }, [metricsData, siteId]);

    useEffect(() => {
      const cleanup = () => {
        Object.values(subscriptionsRef.current).forEach((token) => {
          PubSub.unsubscribe(token);
        });
        subscriptionsRef.current = {};
      };

      cleanup();

      if (loading || !siteId) return;

      const handleUptimeUpdate = (_: any, data: any) => {
        if (data !== null && data !== undefined) {
          const value =
            Array.isArray(data) && data.length > 1
              ? extractMetricValue(data[1])
              : extractMetricValue(data);
          if (value !== null) setUptimeValue(value);
        }
      };

      const handleBatteryUpdate = (_: any, data: any) => {
        if (data !== null && data !== undefined) {
          const value =
            Array.isArray(data) && data.length > 1
              ? extractMetricValue(data[1])
              : extractMetricValue(data);
          if (value !== null) setBatteryValue(value);
        }
      };

      const handleBackhaulUpdate = (_: any, data: any) => {
        if (data !== null && data !== undefined) {
          const value =
            Array.isArray(data) && data.length > 1
              ? extractMetricValue(data[1])
              : extractMetricValue(data);
          if (value !== null) setBackhaulValue(value);
        }
      };

      const uptimeTopic = `stat-${SITE_KPI_TYPES.SITE_UPTIME}-${siteId}`;
      const batteryTopic = `stat-${SITE_KPI_TYPES.BATTERY_CHARGE_PERCENTAGE}-${siteId}`;
      const backhaulTopic = `stat-${SITE_KPI_TYPES.BACKHAUL_SPEED}-${siteId}`;

      subscriptionsRef.current.uptime = PubSub.subscribe(
        uptimeTopic,
        handleUptimeUpdate,
      );
      subscriptionsRef.current.battery = PubSub.subscribe(
        batteryTopic,
        handleBatteryUpdate,
      );
      subscriptionsRef.current.backhaul = PubSub.subscribe(
        backhaulTopic,
        handleBackhaulUpdate,
      );

      return cleanup;
    }, [siteId, loading]);

    const displayAddress = loading
      ? ''
      : truncateText(address || '', maxAddressLength);

    const handleClick = useCallback(
      (event: React.MouseEvent<HTMLButtonElement>) => {
        event.stopPropagation();
        setAnchorEl(event.currentTarget);
      },
      [],
    );

    const handleClose = useCallback(() => {
      setAnchorEl(null);
    }, []);

    const handleMenuClick = useCallback(() => {
      handleSiteNameUpdate(siteId, name);
      handleClose();
    }, [handleSiteNameUpdate, siteId, name, handleClose]);

    const navigateToDetails = useCallback(() => {
      window.location.href = `/console/sites/${siteId}`;
    }, [siteId]);

    const connectionStyles =
      uptimeValue !== null
        ? getStatusStyles('uptime', uptimeValue)
        : {
            icon: <RouterIcon sx={{ color: colors.darkGray }} />,
            color: colors.darkGray,
          };

    const batteryStyles =
      batteryValue !== null
        ? getStatusStyles('battery', batteryValue)
        : {
            icon: <BatteryAlertIcon sx={{ color: colors.darkGray }} />,
            color: colors.darkGray,
          };

    const signalStyles =
      backhaulValue !== null
        ? getStatusStyles('signal', backhaulValue)
        : {
            icon: <SignalCellularOffIcon sx={{ color: colors.darkGray }} />,
            color: colors.darkGray,
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
            {/* TODO: Commenting this out for now as we don't have a way to get subscribers by site yet
            <Box display="flex" alignItems="center" gap={1}>
              {loading || userCount === undefined || userCount === null ? (
                <>
                  <PeopleIcon sx={{ color: colors.darkGray }} />
                  <Typography variant="body2" sx={{ color: colors.darkGray }}>
                    0
                  </Typography>
                </>
              ) : (
                <>
                  <PeopleIcon sx={{ color: colors.darkGray }} />
                  <Typography variant="body2" sx={{ color: colors.darkGray }}>
                    {userCount}
                  </Typography>
                </>
              )}
            </Box>
          */}
            <Box display="flex" alignItems="center" gap={1}>
              {connectionStyles.icon}
              <Typography
                variant="body2"
                sx={{
                  color: connectionStyles.color,
                  display: { xs: 'none', sm: 'block' },
                }}
              >
                {loading || uptimeValue === null
                  ? 'Pending'
                  : uptimeValue <= 0
                    ? 'Offline'
                    : 'Online'}
              </Typography>
            </Box>

            <Box display="flex" alignItems="center" gap={1}>
              {batteryStyles.icon}
              <Typography
                variant="body2"
                sx={{
                  color: batteryStyles.color,
                  display: { xs: 'none', sm: 'block' },
                }}
              >
                {loading || batteryValue === null
                  ? 'Pending'
                  : batteryValue < 20
                    ? 'Critical'
                    : batteryValue < 40
                      ? 'Low'
                      : 'Charged'}
              </Typography>
            </Box>

            <Box display="flex" alignItems="center" gap={1}>
              {signalStyles.icon}
              <Typography
                variant="body2"
                sx={{
                  color: signalStyles.color,
                  display: { xs: 'none', sm: 'block' },
                }}
              >
                {loading || backhaulValue === null
                  ? 'Pending'
                  : backhaulValue < 10
                    ? 'No signal'
                    : backhaulValue < 70
                      ? 'Low signal'
                      : 'Strong'}
              </Typography>
            </Box>
          </Box>
        </CardContent>
      </Card>
    );
  },
  (prevProps, nextProps) => {
    return (
      prevProps.siteId === nextProps.siteId &&
      prevProps.name === nextProps.name &&
      prevProps.address === nextProps.address &&
      prevProps.loading === nextProps.loading &&
      prevProps.userCount === nextProps.userCount &&
      prevProps.metricsData === nextProps.metricsData
    );
  },
);

SiteCard.displayName = 'SiteCard';

export default SiteCard;
