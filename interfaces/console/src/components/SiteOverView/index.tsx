/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Tooltip,
  Stack,
} from '@mui/material';
import React, { useEffect } from 'react';
import { useMemo } from 'react';
import { formatISO } from 'date-fns';

interface SiteOverviewProps {
  siteUptimeSeconds?: number;
  uptimePercentage?: number;
  installationDate?: string;
}
interface DayData {
  date: Date;
  displayDate: string;
  percentage: number;
}
const SiteOverview: React.FC<SiteOverviewProps> = ({
  siteUptimeSeconds = 86400 * 90,
  uptimePercentage = 99,
  installationDate,
}) => {
  const [inputSiteUptimeSeconds, setInputSiteUptimeSeconds] =
    React.useState<number>(siteUptimeSeconds);
  const [inputUptimePercentage, setInputUptimePercentage] =
    React.useState<number>(uptimePercentage);
  const [inputInstallDate, setInputInstallDate] = React.useState<string>(
    installationDate || '',
  );
  const installDate = useMemo(() => {
    if (!inputInstallDate) return new Date();
    const date = new Date(inputInstallDate);
    if (isNaN(date.getTime())) return new Date();
    return new Date(date.getFullYear(), date.getMonth(), date.getDate());
  }, [inputInstallDate]);
  const formatDate = (date: Date): string => {
    return new Intl.DateTimeFormat('en-US', {
      month: 'short',
      day: 'numeric',
    }).format(date);
  };

  const formatFullDate = (date: Date): string => {
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    }).format(date);
  };

  const formatUptime = (seconds: number): string => {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const remainingSeconds = seconds % 60;
    return `${days}d ${hours}h ${minutes}m ${remainingSeconds}s`;
  };
  useEffect(() => {
    setInputSiteUptimeSeconds(siteUptimeSeconds);
    setInputUptimePercentage(uptimePercentage);
    setInputInstallDate(formatISO(installDate));
  }, [siteUptimeSeconds, uptimePercentage, installDate]);
  const generateDaysData = (
    daysCount: number,
    startFromDay: number = 0,
  ): DayData[] => {
    const result: DayData[] = [];
    const startDate = installDate;
    const uptimeDays = Math.ceil(inputSiteUptimeSeconds / 86400);

    for (let i = startFromDay; i < startFromDay + daysCount; i++) {
      const date = new Date(startDate);
      date.setDate(startDate.getDate() + i);

      let dayPercentage = 0;
      if (inputSiteUptimeSeconds > 0 && i < uptimeDays) {
        dayPercentage = inputUptimePercentage;
      }

      result.push({
        date,
        displayDate: formatDate(date),
        percentage: dayPercentage,
      });
    }
    return result;
  };
  const firstThirtyDays = useMemo(() => {
    return generateDaysData(30, 0);
  }, [inputUptimePercentage, inputSiteUptimeSeconds, installDate]);
  const nextThirtyDays = useMemo(() => {
    return generateDaysData(30, 30);
  }, [inputUptimePercentage, inputSiteUptimeSeconds, installDate]);
  const renderBar = (day: DayData, index: number) => {
    const { percentage, displayDate } = day;
    const hasData = percentage > 0;
    const uptimeHeight = (percentage / 100) * 70;
    const isInstallDate = day.date.getTime() === installDate.getTime();

    let tooltipText = `${displayDate}`;

    if (isInstallDate) {
      tooltipText += `\nInstallation Date: ${formatFullDate(installDate)}`;
    }

    if (hasData) {
      tooltipText += `\nUptime: ${percentage.toFixed(1)}%`;
    } else {
      tooltipText += `\nNo uptime data`;
    }

    return (
      <Tooltip key={index} title={tooltipText} placement="top">
        <Box
          sx={{
            height: 70,
            width: 10,
            mx: 0.25,
            borderRadius: 10,
            position: 'relative',
            bgcolor: hasData ? colors.redLight : colors.gray,
            overflow: 'hidden',
            border: isInstallDate ? '2px solid black' : 'none',
          }}
        >
          {hasData && (
            <Box
              sx={{
                position: 'absolute',
                bottom: 0,
                width: '100%',
                height: `${uptimeHeight}px`,
                bgcolor: colors.lightGreen,
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

        <Typography
          variant="subtitle2"
          sx={{
            mb: 4,
            color: uptimePercentage > 98 ? colors.gray : colors.redLight,
          }}
        >
          {uptimePercentage}% uptime over {formatUptime(siteUptimeSeconds)}
        </Typography>

        <Box sx={{ mb: 5 }}>
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
              30 days
            </Typography>
          </Box>
        </Box>

        <Box>
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
              31 days
            </Typography>
            <Typography variant="caption" color="text.secondary">
              60 days
            </Typography>
          </Box>
        </Box>

        <Stack direction="row" spacing={2} sx={{ mt: 4 }}>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Box
              sx={{ width: 10, height: 10, bgcolor: colors.lightGreen, mr: 1 }}
            />
            <Typography variant="caption">Uptime</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Box
              sx={{ width: 10, height: 10, bgcolor: colors.redLight, mr: 1 }}
            />
            <Typography variant="caption">Downtime</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Box sx={{ width: 10, height: 10, bgcolor: colors.gray, mr: 1 }} />
            <Typography variant="caption">No data</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Box
              sx={{
                width: 10,
                height: 10,
                border: '2px solid black',
                mr: 1,
              }}
            />
            <Typography variant="caption">Installation Date</Typography>
          </Box>
        </Stack>
      </CardContent>
    </Card>
  );
};
export default SiteOverview;
