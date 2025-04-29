/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useEffect, useMemo, useState } from 'react';
import {
  Box,
  Typography,
  Tooltip,
  Card,
  CardContent,
  Stack,
  Skeleton,
} from '@mui/material';
import { subDays, addDays, format } from 'date-fns';
import colors from '@/theme/colors';
import { duration } from '@/utils';
import { SiteMetricsStateRes } from '@/client/graphql/generated/subscriptions';
import { SITE_KPI_TYPES } from '@/constants';

interface DayData {
  date: Date;
  displayDate: string;
  percentage: number;
  daysAgo: number | null;
  isInstallDay: boolean;
}

interface SiteOverviewProps {
  installationDate: Date;
  includeFutureDays?: boolean;
  isLoading: boolean;
  siteId: string;
  siteStatMetrics: SiteMetricsStateRes;
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
      const value = uptimePercentageMetric.value;
      const numValue = typeof value === 'number' ? value : parseFloat(value);
      setUptimePercentage(Math.floor(numValue));
    }

    if (uptimeSecondsMetric?.value !== undefined) {
      const value = uptimeSecondsMetric.value;
      const numValue = typeof value === 'number' ? value : parseFloat(value);
      setSiteUptimeSeconds(Math.floor(numValue));
    }
  }, [siteId, siteStatMetrics]);

  useEffect(() => {
    if (!siteId) return;

    const topics = [
      `stat-${SITE_KPI_TYPES.SITE_UPTIME_PERCENTAGE}-${siteId}`,
      `stat-${SITE_KPI_TYPES.SITE_UPTIME}-${siteId}`,
    ];

    const tokens = topics.map((topic) =>
      PubSub.subscribe(topic, (_, value) => {
        if (topic.includes('percentage')) {
          setUptimePercentage(Math.floor(value));
        } else {
          setSiteUptimeSeconds(Math.floor(value));
        }
      }),
    );

    return () => {
      tokens.forEach((token) => PubSub.unsubscribe(token));
    };
  }, [siteId]);
  const isSameDay = (dateA: Date, dateB: Date) => {
    return (
      dateA.getFullYear() === dateB.getFullYear() &&
      dateA.getMonth() === dateB.getMonth() &&
      dateA.getDate() === dateB.getDate()
    );
  };

  const today = useMemo(() => {
    const now = new Date();
    return new Date(now.getFullYear(), now.getMonth(), now.getDate());
  }, []);

  const installDate = useMemo(() => {
    return new Date(
      installationDate.getFullYear(),
      installationDate.getMonth(),
      installationDate.getDate(),
    );
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

    for (let i = 0; i < 60; i++) {
      const date = subDays(today, i);
      const displayDate = format(date, 'MMM d');
      const daysAgo = i;
      const isInstallDay = isSameDay(date, installDate);

      let percentage = 0;

      if (isFutureInstall) {
        percentage = 0;
      } else if (isInstallationToday) {
        percentage = i === 0 ? (uptimePercentage ?? 0) : 0;
      } else {
        const isPastOrEqualToInstall =
          date.getTime() >= installDate.getTime() ||
          isSameDay(date, installDate);
        percentage = isPastOrEqualToInstall ? (uptimePercentage ?? 0) : 0;
      }

      result.push({
        date,
        displayDate,
        percentage,
        daysAgo,
        isInstallDay,
      });
    }

    if (
      includeFutureDays &&
      isInstallationToday &&
      uptimeDays &&
      uptimeDays > 1
    ) {
      const futureDaysToShow = uptimeDays - 1;
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

  const renderBar = (day: DayData, index: number) => {
    const isFutureDay = (day.daysAgo ?? 0) < 0;
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

  return (
    <Card
      sx={{
        borderRadius: 2,
        boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.05)',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <CardContent sx={{ padding: 4, flexGrow: 1 }}>
        <Typography variant="h6" sx={{ mb: 3 }}>
          Site overview
        </Typography>

        {isLoading || uptimePercentage == null ? (
          <Skeleton variant="text" width="60" height={20} sx={{ mb: 2 }} />
        ) : (
          <Typography
            variant="subtitle2"
            sx={{
              mb: 4,
              color:
                (uptimePercentage ?? 0) >= 99
                  ? colors.black70
                  : colors.redLight,
            }}
          >
            {uptimePercentage}% uptime over {duration(siteUptimeSeconds ?? 0)}
          </Typography>
        )}
        <Stack direction="row" spacing={2} alignItems={'center'} sx={{ mb: 2 }}>
          <Stack direction={'row'} alignItems={'center'} spacing={1}>
            <Box
              sx={{
                width: 14,
                height: 14,
                bgcolor: colors.ligthGreen,
              }}
            />
            <Typography variant="caption">Good Uptime</Typography>
          </Stack>
          <Stack direction={'row'} alignItems={'center'} spacing={1}>
            <Box
              sx={{
                width: 14,
                height: 14,
                bgcolor: colors.gray,
              }}
            />
            <Typography variant="caption">No Uptime</Typography>
          </Stack>
          <Stack direction={'row'} alignItems={'center'} spacing={1}>
            <Box
              sx={{
                width: 14,
                height: 14,
                bgcolor: colors.redLight,
              }}
            />
            <Typography variant="caption">Low Uptime</Typography>
          </Stack>

          <Box
            sx={{
              width: 10,
              height: 10,
              bgcolor: colors.gray,
              border: `2px solid ${colors.black70}`,
              mr: 1,
            }}
          />
          <Typography variant="caption">Installation Day</Typography>
        </Stack>
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-end',
            height: 70,
            mb: 2,
            flexDirection: 'row-reverse',
          }}
        >
          {firstThirtyDays.map((day, index) => renderBar(day, index))}
        </Box>
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'space-between',
            mt: 1,
            flexDirection: 'row-reverse',
          }}
        >
          <Typography variant="caption" color="text.secondary">
            Today
          </Typography>
          <Typography variant="caption" color="text.secondary">
            30 days ago
          </Typography>
        </Box>

        <Box
          sx={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-end',
            height: 70,
            mb: 1,
            flexDirection: 'row-reverse',
          }}
        >
          {nextThirtyDays.map((day, index) => renderBar(day, index))}
        </Box>
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'space-between',
            mt: 1,
            flexDirection: 'row-reverse',
          }}
        >
          <Typography variant="caption" color="text.secondary">
            31 days ago
          </Typography>
          <Typography variant="caption" color="text.secondary">
            60 days ago
          </Typography>
        </Box>
      </CardContent>
    </Card>
  );
};

export default SiteOverview;
