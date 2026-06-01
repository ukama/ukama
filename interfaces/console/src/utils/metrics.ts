/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import {
  MetricRes,
  MetricsRes,
  MetricsStateRes,
} from '@/client/graphql/generated/subscriptions';
import { KPI_PLACEHOLDER_VALUE } from '@/constants';
import colors from '@/theme/colors';
import Highcharts, { DashStyleValue } from 'highcharts';

export const getGraphFilterByType = (type: string) => {
  switch (type) {
    case 'DAY':
      return {
        to: Math.round(Date.now() / 1000),
        from: Math.round(Date.now() / 1000) - 86400,
      };
    case 'WEEK':
      return {
        to: Math.round(Date.now() / 1000),
        from: Math.round(Date.now() / 1000) - 604800,
      };
    case 'MONTH':
      return {
        to: Math.round(Date.now() / 1000),
        from: Math.round(Date.now() / 1000) - 2628002,
      };
  }
};

export const extractMetricValue = (value: unknown): number | null => {
  if (Array.isArray(value) && value.length > 1) {
    return typeof value[1] === 'number' ? value[1] : null;
  }
  return typeof value === 'number' ? value : null;
};

export const getMetricValue = (key: string, metrics: MetricsRes) => {
  const metric = metrics.metrics.find((item: MetricRes) => item.type === key);
  return metric?.values ?? [];
};

export const isMetricValue = (key: string, metrics: MetricsRes) => {
  const metric = metrics.metrics.find((item: MetricRes) => item.type === key);
  return (metric && metric.values.length > 1) ?? false;
};

export const getKPIStatValue = (
  id: string,
  loading: boolean,
  statsData: MetricsStateRes | MetricsStateRes,
): string => {
  if (loading || !statsData?.metrics) return KPI_PLACEHOLDER_VALUE;
  const stat = statsData.metrics.find((item) => item.type === id);
  return stat?.value?.toString() ?? KPI_PLACEHOLDER_VALUE;
};

export const findNullZones = (data: [number, number | null][]) => {
  const zones: Highcharts.SeriesZonesOptionsObject[] = [];
  let inNullZone = false;
  let start: number | null = null;

  for (let i = 0; i < data.length; i++) {
    const [x, y] = data[i];

    if (y === null) {
      if (!inNullZone) {
        start = x;
        inNullZone = true;
      }
    } else {
      if (inNullZone && start !== null) {
        zones.push({ value: start });
        zones.push({
          value: data[i - 1][0],
          color: colors.black38,
          dashStyle: 'dash' as DashStyleValue,
        });
        inNullZone = false;
      }
    }
  }

  if (inNullZone && start !== null) {
    zones.push({ value: start });
    zones.push({
      value: data[data.length - 1][0],
      color: colors.black38,
      dashStyle: 'dash' as DashStyleValue,
    });
  }

  return zones;
};

export const generatePlotLines = (
  values: number[] | undefined,
): Highcharts.XAxisPlotLinesOptions[] => {
  if (!values || values.length === 0) {
    return [];
  }

  if (values.length < 3 || values.length > 7) {
    return [];
  }

  return values.slice(1).map((value, index, arr) => ({
    value,
    color:
      index === 0
        ? colors.dullGrey
        : index === arr.length - 2
          ? colors.dullRed
          : index === arr.length - 1
            ? colors.white
            : colors.dullGreen,
    width: 2,
    zIndex: 4,
    dashStyle: 'Dash',
  }));
};

export const formatKPIValue = (
  value: string,
  type: string,
): string | number => {
  switch (type) {
    case 'number':
      return Math.floor(parseFloat(value));
    case 'decimal':
      return parseFloat(value).toFixed(2);
    default:
      return value.toString();
  }
};
