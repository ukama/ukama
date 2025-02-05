/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Box } from '@mui/material';
import { HighchartsReact } from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import PubSub from 'pubsub-js';
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
    return {
      title: {
        text: topic,
      },

      chart: {
        type: 'areaspline',
        events: {
          load: function () {
            const chart: any =
              Highcharts.charts.length > 0
                ? Highcharts.charts.find(
                    (c: any) => c?.title?.textStr === topic,
                  )
                : null;

            if (chart) {
              PubSub.subscribe(topic, (_, data) => {
                if (data && Array.isArray(data) && data.length > 0) {
                  for (let i = 0; i < data.length; i++) {
                    data && chart?.series[0].addPoint(data[i], false, true);
                  }
                  chart.redraw();
                }
              });
            }
          },
        },
      },

      plotOptions: {
        areaspline: {
          color: '#218FF6A2',
          fillColor: {
            linearGradient: { x1: 0, x2: 0, y1: 0, y2: 1 },
            stops: [
              [0, '#218FF66F'],
              [1, '#218FF61B'],
            ],
          },
          threshold: null,
          marker: {
            lineWidth: 1,
            lineColor: null,
            fillColor: 'white',
          },
        },
      },

      time: {
        useUTC: false,
      },

      exporting: {
        enabled: false,
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

      accessibility: {
        enabled: false,
      },

      series: [
        {
          name: title,
          data: (function () {
            const data = [...initData];
            return data;
          })(),
        },
      ],

      xAxis: {
        type: 'datetime',
        title: false,
        labels: {
          enabled: true,
          format: '{value:%H:%M}',
          dateTimeLabelFormats: {
            millisecond: '%H:%M:%S.%L',
            second: '%H:%M:%S',
            minute: '%H:%M',
            hour: '%H:%M',
            day: '%e. %b',
            week: '%e. %b',
            month: "%b '%y",
            year: '%Y',
          },
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
          if (f === 'LIVE') {
            chart.xAxis[0].setExtremes(null, null);
          }
          chart.update(
            {
              navigator: {
                enabled: f === 'ZOOM',
              },
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
