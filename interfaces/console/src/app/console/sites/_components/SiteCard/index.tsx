/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { MetricsStateRes } from '@/client/graphql/generated/subscriptions';
import { SITE_KPI_TYPES } from '@/constants';
import colors from '@/theme/colors';
import { getStatusStyles } from '@/utils';
import BatteryAlertIcon from '@mui/icons-material/BatteryAlert';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import PeopleIcon from '@mui/icons-material/People';
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
  Typography,
} from '@mui/material';
import { useRouter } from 'next/navigation';
import PubSub from 'pubsub-js';
import React, { memo, useCallback, useEffect, useRef, useState } from 'react';
import {
  extractMetricFromPubSubPayload,
  getSiteActiveSubscribers,
  getSiteMetricValue,
  truncateText,
} from './utils';

type PubSubPayload = number | [unknown, number] | unknown[] | [unknown, number | [unknown, number] | unknown[]];
type SubscriptionMap = {
  uptime?: string;
  battery?: string;
  backhaul?: string;
  subscribers?: string;
};

interface SiteCardProps {
  siteId: string;
  name: string;
  address: string;
  userCount?: number;
  loading?: boolean;
  handleSiteNameUpdate: (siteId: string, newSiteName: string) => void;
  maxAddressLength?: number;
  metricsData?: MetricsStateRes;
}


const SiteCard: React.FC<SiteCardProps> = memo(
  ({
    siteId,
    name,
    address,
    handleSiteNameUpdate,
    loading,
    maxAddressLength = 49,
    metricsData,
  }) => {
    const router = useRouter();
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

    const [uptimeValue, setUptimeValue] = useState<number | null>(null);
    const [batteryValue, setBatteryValue] = useState<number | null>(null);
    const [backhaulValue, setBackhaulValue] = useState<number | null>(null);
    const [activeSubscribers, setActiveSubscribers] = useState<number | null>(
      null,
    );

    const subscriptionsRef = useRef<SubscriptionMap>({});

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

      const subscribers = getSiteActiveSubscribers(metricsData, siteId);
      if (subscribers !== null) setActiveSubscribers(subscribers);
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

      const createMetricHandler =
        (setter: React.Dispatch<React.SetStateAction<number | null>>) =>
        (_: string, data: PubSubPayload) => {
          const value = extractMetricFromPubSubPayload(data);
          if (value !== null) setter(value);
        };

      const uptimeTopic = `stat-${SITE_KPI_TYPES.SITE_UPTIME}-${siteId}`;
      const batteryTopic = `stat-${SITE_KPI_TYPES.BATTERY_CHARGE_PERCENTAGE}-${siteId}`;
      const backhaulTopic = `stat-${SITE_KPI_TYPES.BACKHAUL_SPEED}-${siteId}`;
      const subscribersTopic = `stat-${SITE_KPI_TYPES.ACTIVE_SUBSCRIBERS}-${siteId}`;

      subscriptionsRef.current.uptime = PubSub.subscribe(
        uptimeTopic,
        createMetricHandler(setUptimeValue),
      );
      subscriptionsRef.current.battery = PubSub.subscribe(
        batteryTopic,
        createMetricHandler(setBatteryValue),
      );
      subscriptionsRef.current.backhaul = PubSub.subscribe(
        backhaulTopic,
        createMetricHandler(setBackhaulValue),
      );
      subscriptionsRef.current.subscribers = PubSub.subscribe(
        subscribersTopic,
        createMetricHandler(setActiveSubscribers),
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
      // Use Next.js router.push instead of window.location.href so that
      // client-side navigation is preserved (no full-page reload, back button
      // works correctly, prefetch kicks in).
      router.push(`/console/sites/${siteId}`);
    }, [router, siteId]);

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

    const statusTextSx = {
      display: { xs: 'none', sm: 'block' },
    };

    const getSubscribersDisplayValue = () =>
      loading || activeSubscribers === null ? 0 : activeSubscribers;

    const getConnectionLabel = () => {
      if (loading || uptimeValue === null) return 'Pending';
      return uptimeValue <= 0 ? 'Offline' : 'Online';
    };

    const getBatteryLabel = () => {
      if (loading || batteryValue === null) return 'Pending';
      if (batteryValue < 20) return 'Critical';
      if (batteryValue < 40) return 'Low';
      return 'Charged';
    };

    const getSignalLabel = () => {
      if (loading || backhaulValue === null) return 'Pending';
      if (backhaulValue < 10) return 'No signal';
      if (backhaulValue < 70) return 'Low signal';
      return 'Strong';
    };

    const renderMetricItem = (
      icon: React.ReactNode,
      value: string,
      color: string,
    ) => (
      <Box display="flex" alignItems="center" gap={1}>
        {icon}
        <Typography variant="body2" sx={{ color, ...statusTextSx }}>
          {value}
        </Typography>
      </Box>
    );

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
              <PeopleIcon sx={{ color: colors.darkGray }} />
              <Typography variant="body2" sx={{ color: colors.darkGray }}>
                {getSubscribersDisplayValue()}
              </Typography>
            </Box>

            {renderMetricItem(
              connectionStyles.icon,
              getConnectionLabel(),
              connectionStyles.color,
            )}

            {renderMetricItem(
              batteryStyles.icon,
              getBatteryLabel(),
              batteryStyles.color,
            )}

            {renderMetricItem(
              signalStyles.icon,
              getSignalLabel(),
              signalStyles.color,
            )}
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
