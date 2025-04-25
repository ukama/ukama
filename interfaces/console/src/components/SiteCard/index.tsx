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
import React, { useEffect, useState, memo, useCallback, useRef } from 'react';
import PubSub from 'pubsub-js';
import { SiteMetricsStateRes } from '@/client/graphql/generated/subscriptions';
import { SITE_KPI_TYPES } from '@/constants';

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
    (m) => m.type === metricId && m.success == true && m.siteId === siteId,
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
        metricsData &&
        metricsData.metrics &&
        metricsData.metrics.length > 0
      ) {
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
      }
    }, [metricsData, siteId]);

    useEffect(() => {
      Object.values(subscriptionsRef.current).forEach((token) => {
        PubSub.unsubscribe(token);
      });
      subscriptionsRef.current = {};

      const uptimeTopic = `stat-${SITE_KPI_TYPES.SITE_UPTIME}-${siteId}`;
      const batteryChargeTopic = `stat-${SITE_KPI_TYPES.BATTERY_CHARGE_PERCENTAGE}-${siteId}`;
      const backhaulTopic = `stat-${SITE_KPI_TYPES.BACKHAUL_SPEED}-${siteId}`;

      const uptimeToken = PubSub.subscribe(uptimeTopic, (_, data) => {
        if (data && data.length > 1) {
          const value = extractMetricValue(data[1]);
          if (value !== null) setUptimeValue(value);
        }
      });

      const batteryChargeToken = PubSub.subscribe(
        batteryChargeTopic,
        (_, data) => {
          if (data && data.length > 1) {
            const value = extractMetricValue(data[1]);
            if (value !== null) setBatteryValue(value);
          }
        },
      );

      const backhaulToken = PubSub.subscribe(backhaulTopic, (_, data) => {
        if (data && data.length > 1) {
          const value = extractMetricValue(data[1]);
          if (value !== null) setBackhaulValue(value);
        }
      });

      subscriptionsRef.current = {
        uptime: uptimeToken,
        battery: batteryChargeToken,
        backhaul: backhaulToken,
      };

      return () => {
        Object.values(subscriptionsRef.current).forEach((token) => {
          PubSub.unsubscribe(token);
        });
        subscriptionsRef.current = {};
      };
    }, [siteId]);

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
