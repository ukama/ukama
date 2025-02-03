/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useMemo } from 'react';
import { Box, Stack, Typography, useTheme } from '@mui/material';
import Highcharts from 'highcharts/highstock';
import HighchartsReact from 'highcharts-react-official';
import { MetricsRes } from '@/client/graphql/generated/subscriptions';
import { getMetricValue } from '@/utils';

interface SiteOverviewProps {
  metrics: MetricsRes;
  loading: boolean;
}

const SiteOverview: React.FC<SiteOverviewProps> = ({ metrics, loading }) => {
  const theme = useTheme();

  const currentValues = useMemo(() => {
    const powerValue = getMetricValue('solar_panel_power', metrics);
    const batteryVoltage = getMetricValue('battery_voltage', metrics);
    const panelVoltage = getMetricValue('solar_panel_voltage', metrics);

    return {
      power:
        Array.isArray(powerValue) && powerValue.length > 0
          ? (powerValue[powerValue.length - 1]?.[1] ?? 0)
          : 0,
      batteryVoltage:
        Array.isArray(batteryVoltage) && batteryVoltage.length > 0
          ? (batteryVoltage[batteryVoltage.length - 1]?.[1] ?? 0)
          : 0,
      solarVoltage:
        Array.isArray(panelVoltage) && panelVoltage.length > 0
          ? (panelVoltage[panelVoltage.length - 1]?.[1] ?? 0)
          : 0,
    };
  }, [metrics]);

  const chartData = useMemo(() => {
    const powerData = getMetricValue('solar_panel_power', metrics);
    const batteryData = getMetricValue('battery_voltage', metrics);
    const solarVoltageData = getMetricValue('solar_panel_voltage', metrics);

    return {
      power: Array.isArray(powerData)
        ? powerData.map(([time, value]) => [time * 1000, Number(value) || 0])
        : [],
      battery: Array.isArray(batteryData)
        ? batteryData.map(([time, value]) => [time * 1000, Number(value) || 0])
        : [],
      solarVoltage: Array.isArray(solarVoltageData)
        ? solarVoltageData.map(([time, value]) => [
            time * 1000,
            Number(value) || 0,
          ])
        : [],
    };
  }, [metrics]);

  const chartOptions: Highcharts.Options = {
    chart: {
      height: 160,
      style: {
        fontFamily: theme.typography.fontFamily,
      },
    },
    title: {
      text: 'Power & Battery Overview',
      style: {
        fontSize: '14px',
      },
    },
    xAxis: {
      type: 'datetime',
      labels: {
        format: '{value:%H:%M}',
      },
    },
    yAxis: [
      {
        title: {
          text: 'Power (W)',
          style: {
            color: theme.palette.primary.main,
          },
        },
        labels: {
          style: {
            color: theme.palette.primary.main,
          },
        },
      },
      {
        title: {
          text: 'Voltage (V)',
          style: {
            color: theme.palette.success.main,
          },
        },
        opposite: true,
        labels: {
          style: {
            color: theme.palette.success.main,
          },
        },
      },
    ],
    series: [
      {
        name: 'Solar Power',
        type: 'line',
        data: chartData.power,
        color: theme.palette.primary.main,
        yAxis: 0,
        tooltip: {
          valueSuffix: ' W',
        },
      },
      {
        name: 'Battery Voltage',
        type: 'line',
        data: chartData.battery,
        color: theme.palette.success.main,
        yAxis: 1,
        tooltip: {
          valueSuffix: ' V',
        },
      },
      {
        name: 'Solar Panel Voltage',
        type: 'line',
        data: chartData.solarVoltage,
        color: theme.palette.warning.main,
        yAxis: 1,
        tooltip: {
          valueSuffix: ' V',
        },
      },
    ],
    legend: {
      enabled: true,
      align: 'left',
      verticalAlign: 'top',
    },
    tooltip: {
      shared: true,
    },
    credits: {
      enabled: false,
    },
    time: {
      useUTC: false,
    },
  };

  const MetricIndicator: React.FC<{
    color: string;
    label: string;
    value: string | number;
    unit: string;
  }> = ({ color, label, value, unit }) => (
    <Stack direction="row" alignItems="center" spacing={1}>
      <Box
        sx={{
          width: 10,
          height: 10,
          borderRadius: '50%',
          bgcolor: color,
        }}
      />
      <Typography variant="body2">
        {label}:{' '}
        {typeof value === 'number'
          ? `${Number(value).toFixed(1)}${unit}`
          : value}
      </Typography>
    </Stack>
  );

  return (
    <Box sx={{ p: 2 }}>
      <Typography variant="h6" sx={{ mb: 2 }}>
        Site Overview
      </Typography>

      <Stack direction="row" spacing={3} sx={{ mb: 3 }}>
        <MetricIndicator
          color={theme.palette.primary.main}
          label="Solar Power"
          value={currentValues.power}
          unit="W"
        />
        <MetricIndicator
          color={theme.palette.success.main}
          label="Battery Voltage"
          value={currentValues.batteryVoltage}
          unit="V"
        />
        <MetricIndicator
          color={theme.palette.warning.main}
          label="Solar Voltage"
          value={currentValues.solarVoltage}
          unit="V"
        />
      </Stack>

      <Box sx={{ height: 160, width: '100%' }}>
        <HighchartsReact highcharts={Highcharts} options={chartOptions} />
      </Box>
    </Box>
  );
};

export default SiteOverview;
