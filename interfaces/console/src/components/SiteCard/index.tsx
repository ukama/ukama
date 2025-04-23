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

interface SiteMetricsState {
  site_uptime_seconds: number | null;
  battery_charge_percentage: number | null;
  backhaul_speed: number | null;
  [key: string]: number | null;
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
    const [metrics, setMetrics] = useState<SiteMetricsState>({
      site_uptime_seconds: null,
      battery_charge_percentage: null,
      backhaul_speed: null,
    });
    const updateCount = useRef(0);
    const hasLiveData = useRef<Record<string, boolean>>({});

    useEffect(() => {
      if (
        metricsData &&
        metricsData.metrics &&
        metricsData.metrics.length > 0
      ) {
        const siteMetrics = metricsData.metrics.filter(
          (metric) => metric.siteId === siteId,
        );
        if (siteMetrics.length > 0) {
          setMetrics((prev) => {
            const newMetrics = { ...prev };
            siteMetrics.forEach((metric) => {
              const metricType = metric.type;
              const metricValue = extractMetricValue(metric.value);
              if (!hasLiveData.current[metricType] && metricValue !== null) {
                newMetrics[metricType] = metricValue;
              }
            });
            updateCount.current += 1;
            return newMetrics;
          });
        }
      }
    }, [metricsData, siteId]);

    useEffect(() => {
      const siteMetricsTopic = `site-metrics-${siteId}`;
      const token = PubSub.subscribe(
        siteMetricsTopic,
        (_, data: MetricsUpdateData) => {
          if (data && typeof data === 'object') {
            if (
              'metrics' in data &&
              data.metrics &&
              Array.isArray(data.metrics) &&
              data.metrics.length > 0
            ) {
              setMetrics((prev) => {
                const newMetrics = { ...prev };
                data.metrics!.forEach((metric) => {
                  const metricValue = extractMetricValue(metric.value);
                  newMetrics[metric.type] = metricValue;
                  hasLiveData.current[metric.type] = true;
                });
                updateCount.current += 1;
                return newMetrics;
              });
            } else if (
              'type' in data &&
              data.type &&
              'value' in data &&
              data.value !== undefined
            ) {
              const metricType = data.type;
              const metricValue = extractMetricValue(data.value);
              setMetrics((prev) => {
                const newState = { ...prev };
                newState[metricType] = metricValue;
                hasLiveData.current[metricType] = true;
                updateCount.current += 1;
                return newState;
              });
            }
          }
        },
      );

      const statTopics = [
        'site_uptime_seconds',
        'battery_charge_percentage',
        'backhaul_speed',
      ];
      const statTokens = statTopics.map((metricType) => {
        const topic = `stat-${metricType}`;
        return PubSub.subscribe(topic, (_, data) => {
          if (!data || typeof data !== 'object') return;
          if ('siteId' in data && data.siteId && data.siteId !== siteId) return;
          let metricValue;
          if ('value' in data) {
            metricValue = extractMetricValue(data.value);
          } else {
            metricValue = extractMetricValue(data);
          }
          setMetrics((prev) => {
            const newState = { ...prev };
            newState[metricType] = metricValue;
            hasLiveData.current[metricType] = true;
            updateCount.current += 1;
            return newState;
          });
        });
      });

      return () => {
        PubSub.unsubscribe(token);
        statTokens.forEach((token) => PubSub.unsubscribe(token));
        hasLiveData.current = {};
      };
    }, [siteId]);

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

    const uptimeValue = getSiteMetricValue(
      'site_uptime_seconds',
      metricsData,
      siteId,
    );
    const batteryValue = getSiteMetricValue(
      'battery_charge_percentage',
      metricsData,
      siteId,
    );
    const backhaulValue = getSiteMetricValue(
      'backhaul_speed',
      metricsData,
      siteId,
    );

    const connectionValue =
      metrics.site_uptime_seconds !== null
        ? metrics.site_uptime_seconds
        : uptimeValue;

    const batteryLevel =
      metrics.battery_charge_percentage !== null
        ? metrics.battery_charge_percentage
        : batteryValue;

    const signalLevel =
      metrics.backhaul_speed !== null ? metrics.backhaul_speed : backhaulValue;

    const connectionStyles =
      connectionValue !== null
        ? getStatusStyles('uptime', connectionValue)
        : { icon: null, color: colors.darkGray };

    const batteryStyles =
      batteryLevel !== null
        ? getStatusStyles('battery', batteryLevel)
        : { icon: null, color: colors.darkGray };

    const signalStyles =
      signalLevel !== null
        ? getStatusStyles('signal', signalLevel)
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
            {/* User Count */}
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

            {/* Connection Status */}
            <Box display="flex" alignItems="center" gap={1}>
              {loading || connectionValue === null ? (
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
                    {connectionValue <= 0 ? 'Offline' : 'Online'}
                  </Typography>
                </>
              )}
            </Box>

            {/* Battery Status */}
            <Box display="flex" alignItems="center" gap={1}>
              {loading || batteryLevel === null ? (
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
                    {batteryLevel < 20
                      ? 'Critical'
                      : batteryLevel < 40
                        ? 'Low'
                        : 'Charged'}
                  </Typography>
                </>
              )}
            </Box>

            {/* Signal Status */}
            <Box display="flex" alignItems="center" gap={1}>
              {loading || signalLevel === null ? (
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
                    {signalLevel < 10
                      ? 'No signal'
                      : signalLevel < 70
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
