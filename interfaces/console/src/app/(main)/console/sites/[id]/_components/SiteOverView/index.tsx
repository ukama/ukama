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
import { duration } from '@/utils';
import {
  Box,
  Card,
  CardContent,
  Skeleton,
  Stack,
  Tooltip,
  Typography,
} from '@mui/material';
import { addDays, format, isSameDay, startOfDay, subDays } from 'date-fns';
import React, { useEffect, useMemo, useState } from 'react';

interface DayData {
  date: Date;
  displayDate: string;
  percentage: number;
  daysAgo: number;
  isInstallDay: boolean;
}

interface SiteOverviewProps {
  installationDate: Date;
  includeFutureDays?: boolean;
  isLoading: boolean;
  siteId: string;
  siteStatMetrics: MetricsStateRes;
}

const SiteOverview: React.FC<SiteOverviewProps> = ({
  installationDate,
  includeFutureDays = true,
  isLoading,
  siteId,
  siteStatMetrics,
}) => {
  const [uptimePercentage, setUptimePercentage] = useState<number | null>(null);
  const [siteUptimeSeconds, setSiteUptimeSeconds] = useState<number | null>(
    null,
  );

  const parseMetricNumber = (value: unknown) => {
    if (typeof value === 'number') return value;
    if (typeof value === 'string') return parseFloat(value);
    return Number.NaN;
  };

  useEffect(() => {
    if (!siteId || !siteStatMetrics?.metrics?.length) return;
    const siteMetrics = siteStatMetrics.metrics.filter(
      (metric) => metric.siteId === siteId && metric.success,
    );

    const uptimePercentageMetric = siteMetrics.find(
      (metric) => metric.type === SITE_KPI_TYPES.SITE_UPTIME_PERCENTAGE,
    );

    const uptimeSecondsMetric = siteMetrics.find(
      (metric) => metric.type === SITE_KPI_TYPES.SITE_UPTIME,
    );
    if (uptimePercentageMetric?.value !== undefined) {
      const numValue = parseMetricNumber(uptimePercentageMetric.value);
      // eslint-disable-next-line react-hooks/set-state-in-effect
      if (!Number.isNaN(numValue)) setUptimePercentage(Math.round(numValue));
    }

    if (uptimeSecondsMetric?.value !== undefined) {
      const numValue = parseMetricNumber(uptimeSecondsMetric.value);
       
      if (!Number.isNaN(numValue)) setSiteUptimeSeconds(Math.floor(numValue));
    }
  }, [siteId, siteStatMetrics]);

  useEffect(() => {
    if (!siteId) return;

    const tokens = [
      PubSub.subscribe(
        `stat-${SITE_KPI_TYPES.SITE_UPTIME_PERCENTAGE}`,
        (topic, value) => {
          if (value.length > 0) {
            setUptimePercentage(Math.round(value[1]));
          }
        },
      ),

      PubSub.subscribe(`stat-${SITE_KPI_TYPES.SITE_UPTIME}`, (topic, value) => {
        if (value.length > 0) {
          setSiteUptimeSeconds(Math.floor(value[1]));
        }
      }),
    ];

    return () => {
      tokens.forEach((token) => PubSub.unsubscribe(token));
    };
  }, [siteId]);

  const today = useMemo(() => {
    return startOfDay(new Date());
  }, []);

  const installDate = useMemo(() => {
    return startOfDay(installationDate);
  }, [installationDate]);

  const isInstallationToday = useMemo(
    () => isSameDay(today, installDate),
    [today, installDate],
  );
  const isFutureInstall = useMemo(
    () => installDate.getTime() > today.getTime(),
    [today, installDate],
  );

  const uptimeDays = useMemo(() => {
    return siteUptimeSeconds ? siteUptimeSeconds / 86400 : null;
  }, [siteUptimeSeconds]);

  const generateDaysData = useMemo(() => {
    const result: DayData[] = [];
    const percentageForDate = (date: Date, daysAgo: number) => {
      if (isFutureInstall) return 0;
      if (isInstallationToday)
        return daysAgo === 0 ? (uptimePercentage ?? 0) : 0;
      return date.getTime() >= installDate.getTime()
        ? (uptimePercentage ?? 0)
        : 0;
    };

    for (let i = 0; i < 60; i++) {
      const date = subDays(today, i);
      const displayDate = format(date, 'MMM d');
      const daysAgo = i;
      const isInstallDay = isSameDay(date, installDate);

      result.push({
        date,
        displayDate,
        percentage: percentageForDate(date, daysAgo),
        daysAgo,
        isInstallDay,
      });
    }

    if (
      includeFutureDays &&
      isInstallationToday &&
      uptimeDays != null &&
      uptimeDays > 1
    ) {
      const futureDaysToShow = Math.min(
        60,
        Math.max(0, Math.floor(uptimeDays) - 1),
      );
      if (futureDaysToShow > 0)
        result.splice(60 - futureDaysToShow, futureDaysToShow);

      for (let i = 1; i <= futureDaysToShow; i++) {
        const date = addDays(today, i);
        const displayDate = format(date, 'MMM d');
        const daysAgo = -i;

        result.unshift({
          date,
          displayDate,
          percentage: uptimePercentage ?? 0,
          daysAgo,
          isInstallDay: false,
        });
      }
    }

    return result;
  }, [
    today,
    installDate,
    uptimePercentage,
    isInstallationToday,
    isFutureInstall,
    uptimeDays,
    includeFutureDays,
  ]);

  const firstThirtyDays = generateDaysData.slice(0, 30);
  const nextThirtyDays = generateDaysData.slice(30, 60);

  const legendItems: Array<{
    label: string;
    boxSx: Record<string, unknown>;
  }> = [
    {
      label: 'Good Uptime',
      boxSx: { width: 14, height: 14, bgcolor: colors.ligthGreen },
    },
    {
      label: 'No Uptime',
      boxSx: { width: 14, height: 14, bgcolor: colors.gray },
    },
    {
      label: 'Low Uptime',
      boxSx: { width: 14, height: 14, bgcolor: colors.redLight },
    },
    {
      label: 'Installation Day',
      boxSx: {
        width: 10,
        height: 10,
        bgcolor: colors.gray,
        border: `2px solid ${colors.black70}`,
      },
    },
  ];

  const renderBar = (day: DayData, index: number) => {
    const isFutureDay = day.daysAgo < 0;
    let tooltipContent = `${day.displayDate}: ${day.percentage}% uptime`;
    if (day.isInstallDay) {
      tooltipContent += ' (Installation day)';
    } else if (isFutureDay) {
      tooltipContent += ' (Projected)';
    }
    const uptimeHeight = (day.percentage / 100) * 70;

    return (
      <Tooltip key={index} title={tooltipContent} placement="top">
        <Box
          sx={{
            height: 70,
            width: 10,
            mx: 0.25,
            borderRadius: 10,
            position: 'relative',
            bgcolor: day.percentage > 0 ? colors.redLight : colors.gray,
            overflow: 'hidden',
            border: day.isInstallDay ? `2px solid ${colors.black70}` : 'none',
          }}
        >
          {day.percentage > 0 && (
            <Box
              sx={{
                position: 'absolute',
                bottom: 0,
                width: '100%',
                height: `${uptimeHeight}px`,
                bgcolor: colors.ligthGreen,
              }}
            />
          )}
        </Box>
      </Tooltip>
    );
  };

  const renderBarsSection = (
    days: DayData[],
    leftLabel: string,
    rightLabel: string,
  ) => {
    return (
      <Stack direction="column" spacing={0.5}>
        <Box
          sx={{
            height: 70,
            display: 'flex',
            alignItems: 'flex-end',
            justifyContent: 'space-between',
            flexDirection: 'row-reverse',
          }}
        >
          {days.map((day, index) => renderBar(day, index))}
        </Box>
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            flexDirection: 'row-reverse',
          }}
        >
          <Typography variant="caption" color="text.secondary">
            {leftLabel}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            {rightLabel}
          </Typography>
        </Box>
      </Stack>
    );
  };

  return (
    <Card
      sx={{
        borderRadius: 2,
        boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.05)',
        height: '100%',
        width: '100%',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <CardContent sx={{ padding: 2, flexGrow: 1, minHeight: 0 }}>
        <Stack direction="column" spacing={2}>
          <Typography variant="h6">Site overview</Typography>

          {isLoading || uptimePercentage == null ? (
            <Skeleton variant="text" width="60" height={20} />
          ) : (
            <Typography
              variant="subtitle2"
              sx={{
                color:
                  (uptimePercentage ?? 0) >= 99
                    ? colors.black70
                    : colors.redLight,
              }}
            >
              {uptimePercentage}% uptime over {duration(siteUptimeSeconds ?? 0)}
            </Typography>
          )}

          <Stack direction="row" spacing={2} alignItems="center">
            {legendItems.map((item) => (
              <Stack
                key={item.label}
                direction="row"
                alignItems="center"
                spacing={1}
              >
                <Box sx={item.boxSx} />
                <Typography variant="caption">{item.label}</Typography>
              </Stack>
            ))}
          </Stack>
          {renderBarsSection(firstThirtyDays, 'Today', '30 days ago')}
          {renderBarsSection(nextThirtyDays, '31 days ago', '60 days ago')}
        </Stack>
      </CardContent>
    </Card>
  );
};

export default SiteOverview;
