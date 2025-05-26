/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { findNullZones, formatKPIValue, generatePlotLines } from '@/utils';
import { Box } from '@mui/material';
import HighchartsReact from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import { forwardRef, useMemo, useRef, useState } from 'react';
import GraphTitleWrapper from '../GraphTitleWrapper';
import './linechart.css';

const initDataFixes = (data: any) => {
  return data.map((point: any) => {
    // If the y-value is -1, it indicates missing or invalid data, so we nullify it.
    let y = point[1];
    if (point.length > 0 && y === -1) {
      y = null;
    }
    return [point[0], y];
  });
};
interface ChartProps {
  options: Highcharts.Options;
}

const Chart = forwardRef<HighchartsReact.RefObject, ChartProps>(
  ({ options }, ref) => (
    <HighchartsReact
      highcharts={Highcharts}
      constructorType="stockChart"
      options={options}
      ref={ref}
    />
  ),
);
Chart.displayName = 'Chart';
interface ILineChart {
  topic: string;
  initData: any;
  title?: string;
  from: number;
  hasData?: boolean;
  loading?: boolean;
  yunit?: string;
  format: string;
  tickInterval?: number;
  tickPositions?: number[];
}

const LineChart = ({
  topic,
  yunit,
  format,
  hasData,
  initData,
  title = '',
  loading = false,
  tickInterval = undefined,
  tickPositions = undefined,
}: ILineChart) => {
  const [navigatorEnabled, setNavigatorEnabled] = useState<boolean>(false);
  const chartRef = useRef<HighchartsReact.RefObject>(null);

  const fixedInitData = useMemo(() => initDataFixes(initData), [initData]);

  const chartOptions = useMemo<Highcharts.Options>(
    () => ({
      title: {
        text: topic,
      },

      chart: {
        type: 'spline',
        zooming: {
          mouseWheel: false,
        },

        events: {
          load: function () {
            PubSub.subscribe(`stat-${topic}`, (_, data) => {
              const chart: any =
                Highcharts.charts.length > 0
                  ? Highcharts.charts.find(
                    (c: any) => c?.title?.textStr === topic,
                  )
                  : null;
              if (chart && data.length > 0) {
                const series = chart.series[0];
                series.addPoint(
                  [data[0], formatKPIValue(data[1], format)],
                  true,
                  true,
                  true,
                );
              }
            });
          },
        },
      },

      time: {
        timezone: undefined,
      },

      plotOptions: {
        series: {
          color: colors.primaryMain,
        },
      },

      navigator: {
        enabled: true,
        maskFill: 'rgba(33, 144, 246, 0.15)',
        handles: {
          symbols: ['doublearrow', 'doublearrow'],
          lineWidth: 1,
          width: 9,
          height: 17,
        },
        xAxis: {
          labels: {
            format: '{value:%H:%M}',
          },
        },
      },

      scrollbar: {
        enabled: false,
      },

      xAxis: {
        type: 'datetime',
        tickAmount: 6,
        tickInterval: 1000 * 60 * 31,
        labels: {
          enabled: true,
          format: '{value:%H:%M}',
        },
      },
      yAxis: {
        endOnTick: true,
        max: tickPositions
          ? tickPositions[tickPositions.length - 1]
          : undefined,
        min: tickPositions ? tickPositions[0] : undefined,
        opposite: false,
        gridLineDashStyle: 'Dash',
        tickPositions: tickPositions,
        gridLineWidth: tickPositions ? 0 : 2,
        tickAmount: tickPositions?.length ?? 5,
        tickInterval: tickInterval,
        labels: {
          y: 5,
          formatter: function (v: any) {
            return `${v.value}${yunit}`;
          },
        },

        plotLines: [...generatePlotLines(tickPositions)],
      },

      series: [
        {
          name: title,
          zoneAxis: 'x',
          type: 'spline',
          connectNulls: true,
          data: fixedInitData,
          zones: findNullZones(fixedInitData),
          tooltip: {
            valueDecimals: format === 'number' ? 0 : 2,
          },
        },
      ],
    }),
    [fixedInitData, tickInterval, tickPositions, title, topic, yunit],
  );

  return (
    <GraphTitleWrapper
      title={title}
      hasData={hasData}
      variant="subtitle1"
      handleFilterChange={(f: string) => {
        setNavigatorEnabled(f !== 'LIVE');
        if (f === 'LIVE' && chartRef.current && chartRef.current.chart) {
          chartRef.current.chart.xAxis[0].setExtremes(
            undefined,
            undefined,
            true,
          );
        }
      }}
      loading={loading ?? !initData}
    >
      <Box sx={{ width: '100%' }}>
        <div
          className={`chart-container ${navigatorEnabled ? '' : 'hide-navigator'}`}
        >
          <Chart options={chartOptions} ref={chartRef} />
        </div>
      </Box>
    </GraphTitleWrapper>
  );
};

export default LineChart;
