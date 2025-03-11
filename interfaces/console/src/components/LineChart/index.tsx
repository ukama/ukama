/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { METRIC_RANGE_10800 } from '@/constants';
import { colors } from '@/theme';
import { findNullZones } from '@/utils';
import { Box } from '@mui/material';
import HighchartsReact from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import GraphTitleWrapper from '../GraphTitleWrapper';

interface ILineChart {
  topic: string;
  initData: any;
  title?: string;
  from: number;
  hasData?: boolean;
  loading?: boolean;
}

const LineChart = ({
  topic,
  hasData,
  initData,
  title = '',
  loading = false,
  from: metricFrom,
}: ILineChart) => {
  const getOptions = (topic: string, title: string, initData: any) => {
    const data: any = [];
    if (Array.isArray(initData)) {
      initData.forEach((point: any) => {
        let y = point[1];
        if (point.length > 0 && y === 0) {
          y = null;
        }
        data.push([point[0], y]);
      });
    }

    return {
      title: {
        text: topic,
      },

      chart: {
        type: 'spline',
        events: {
          load: function () {
            console.log(topic);
            PubSub.subscribe(topic, (_, data) => {
              const chart: any =
                Highcharts.charts.length > 0
                  ? Highcharts.charts.find(
                      (c: any) => c?.title?.textStr === topic,
                    )
                  : null;
              if (chart && data.length > 0) {
                const series = chart.series[0];
                series.addPoint(
                  data,
                  true,
                  series.data.length > METRIC_RANGE_10800,
                  true,
                );
              }
            });
          },
        },
      },

      plotOptions: {
        series: {
          color: colors.primaryMain,
        },
      },

      navigator: {
        enabled: false,
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

      time: {
        useUTC: false,
      },

      series: [
        {
          name: title,
          data: (function () {
            return data;
          })(),
          connectNulls: true,
          zoneAxis: 'x',
          zones: findNullZones(data),
        },
      ],

      xAxis: {
        type: 'datetime',
        title: false,
        tickInterval: 1000 * 60 * 30,
        labels: {
          enabled: true,
          format: '{value:%H:%M}',
        },
      },

      yAxis: {
        title: false,
      },
    };
  };

  return (
    <GraphTitleWrapper
      title={title}
      hasData={hasData}
      variant="subtitle1"
      handleFilterChange={(f: string) => {
        const chart: any =
          Highcharts.charts.length > 0
            ? Highcharts.charts.find((c: any) => c?.title?.textStr === topic)
            : null;

        if (chart) {
          const series = chart.series[0].data.map((point: any) => {
            return [point.x, point.y];
          });
          if (f === 'LIVE') {
            chart.xAxis[0].setExtremes(null, null);
          }
          chart.update(
            {
              navigator: {
                enabled: f === 'ZOOM',
              },

              series: [
                {
                  name: title,
                  data: (function () {
                    const data = [...series];
                    return data;
                  })(),
                },
              ],
            },
            true,
          );
        }
      }}
      loading={loading ?? !initData}
    >
      <Box sx={{ width: '100%' }}>
        <HighchartsReact
          key={topic}
          highcharts={Highcharts}
          options={getOptions(topic, title, initData)}
        />
      </Box>
    </GraphTitleWrapper>
  );
};

export default LineChart;
