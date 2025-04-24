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
import { SiteMetricsStateRes } from '@/client/graphql/generated/subscriptions';

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

interface SiteMetric {
  type: string;
  value: number | [number, number];
}

interface MetricsUpdateData {
  metrics?: SiteMetric[] | null;
  type?: string;
  value?: number | [number, number];
  siteId?: string;
  nodeId?: string;
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
    (m) => m.type === metricId && m.siteId === siteId,
  );

  return metric ? extractMetricValue(metric.value) : null;
};

const extractMetricValue = (value: any): number | null => {
  if (Array.isArray(value) && value.length > 1) {
    return typeof value[1] === 'number' ? value[1] : null;
  }
  return typeof value === 'number' ? value : null;
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

    useEffect(() => {
      if (
        metricsData &&
        metricsData.metrics &&
        metricsData.metrics.length > 0
      ) {
        const uptime = getSiteMetricValue(
          'site_uptime_seconds',
          metricsData,
          siteId,
        );
        if (uptime !== null) setUptimeValue(uptime);

        const batteryCharge = getSiteMetricValue(
          'battery_charge_percentage',
          metricsData,
          siteId,
        );

        const backhaul = getSiteMetricValue(
          'backhaul_speed',
          metricsData,
          siteId,
        );
        if (backhaul !== null) setBackhaulValue(backhaul);
      }
    }, [metricsData, siteId]);

    useEffect(() => {
      const uptimeToken = PubSub.subscribe(
        `stat-site_uptime_seconds-${siteId}`,
        (_, data) => {
          if (data && data.length > 1) {
            const value = extractMetricValue(data[1]);
            if (value !== null) setUptimeValue(value);
          }
        },
      );

      const batteryChargeToken = PubSub.subscribe(
        `stat-battery_charge_percentage-${siteId}`,
        (_, data) => {
          if (data && data.length > 1) {
            const value = extractMetricValue(data[1]);
            if (value !== null) setBatteryValue(value);
          }
        },
      );

      const backhaulToken = PubSub.subscribe(
        `stat-backhaul_speed-${siteId}`,
        (_, data) => {
          if (data && data.length > 1) {
            const value = extractMetricValue(data[1]);
            if (value !== null) setBackhaulValue(value);
          }
        },
      );

      return () => {
        PubSub.unsubscribe(uptimeToken);
        PubSub.unsubscribe(batteryChargeToken);
        PubSub.unsubscribe(backhaulToken);
      };
    }, [siteId, batteryValue]);

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

    const connectionStyles =
      uptimeValue !== null
        ? getStatusStyles('uptime', uptimeValue)
        : { icon: null, color: colors.darkGray };

    const batteryStyles =
      batteryValue !== null
        ? getStatusStyles('battery', batteryValue)
        : { icon: null, color: colors.darkGray };

    const signalStyles =
      backhaulValue !== null
        ? getStatusStyles('signal', backhaulValue)
        : { icon: null, color: colors.darkGray };

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
              {loading || uptimeValue === null ? (
                <>
                  <Skeleton width={24} height={24} variant="circular" />
                  <Skeleton width={60} />
                </>
              ) : (
                <>
                  {connectionStyles.icon}
                  <Typography
                    variant="body2"
                    sx={{
                      color: connectionStyles.color,
                      display: { xs: 'none', sm: 'block' },
                    }}
                  >
                    {uptimeValue <= 0 ? 'Offline' : 'Online'}
                  </Typography>
                </>
              )}
            </Box>

            <Box display="flex" alignItems="center" gap={1}>
              {loading || batteryValue === null ? (
                <>
                  <Skeleton width={24} height={24} variant="circular" />
                  <Skeleton width={70} />
                </>
              ) : (
                <>
                  {batteryStyles.icon}
                  <Typography
                    variant="body2"
                    sx={{
                      color: batteryStyles.color,
                      display: { xs: 'none', sm: 'block' },
                    }}
                  >
                    {batteryValue < 20
                      ? 'Critical'
                      : batteryValue < 40
                        ? 'Low'
                        : 'Charged'}
                  </Typography>
                </>
              )}
            </Box>

            <Box display="flex" alignItems="center" gap={1}>
              {loading || backhaulValue === null ? (
                <>
                  <Skeleton width={24} height={24} variant="circular" />
                  <Skeleton width={60} />
                </>
              ) : (
                <>
                  {signalStyles.icon}
                  <Typography
                    variant="body2"
                    sx={{
                      color: signalStyles.color,
                      display: { xs: 'none', sm: 'block' },
                    }}
                  >
                    {backhaulValue < 10
                      ? 'No signal'
                      : backhaulValue < 70
                        ? 'Low signal'
                        : 'Strong'}
                  </Typography>
                </>
              )}
            </Box>
          </Box>
        </CardContent>
      </Card>
    );
  },
);

SiteCard.displayName = 'SiteCard';

export default SiteCard;
