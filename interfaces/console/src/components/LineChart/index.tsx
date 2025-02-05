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
import { useState } from 'react';
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
  const [filter, setFilter] = useState<string>('LIVE');
  const [rangeValues, setRangeValues] = useState({
    to: initData.length > 0 ? initData[initData.length - 1][0] : 0,
    from:
      initData.length > 0
        ? initData[initData.length - 1][0]
        : 0 - 15 * 60 * 1000,
    range:
      initData.length > 0
        ? initData[initData.length - 1][0]
        : 0 -
          (initData.length > 0
            ? initData[initData.length - 1][0]
            : 0 - 15 * 60 * 1000),
  });

  const getOptions = (topic: string, title: string, initData: any) => {
    let count = 0;
    return {
      title: {
        text: topic,
        align: 'left',
      },
      chart: {
        type: 'spline',
        scrollablePlotArea: {
          minWidth: 600,
          scrollPositionX: 1,
        },
        legend: { enabled: false },
        events: {
          load: function () {
            const chart: any =
              Highcharts.charts.length > 0
                ? Highcharts.charts.find(
                    (c: any) => c?.title?.textStr === topic,
                  )
                : null;

            if (chart) {
              const series: any = chart?.series[0];
              PubSub.subscribe(topic, (_, data) => {
                if (topic === chart?.title?.textStr && series) {
                  if (count === 30) {
                    count = 0;
                    chart.xAxis[0].setExtremes(
                      data[0] - rangeValues.range,
                      data[0],
                    );
                    series.addPoint(data, true, true);
                  } else {
                    count++;
                    series.addPoint(data, false, false);
                  }
                }
              });
            }
          },
        },
      },

      time: {
        useUTC: false,
      },

      exporting: {
        enabled: true,
      },

      navigator: {
        enabled: true,
        maskFill: 'rgba(33, 144, 246, 0.15)',
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
        },
      },
      yAxis: {
        title: false,
      },
    };
  };

  const setRange = (f: string) => {
    if (filter === f) return;

    const chart: any =
      Highcharts.charts.length > 0
        ? Highcharts.charts.find((c: any) => c?.title?.textStr === topic)
        : null;

    const series = chart.series[0];
    const data = series.data;

    if (chart) {
      const now = data[data.length - 1].x;
      const from15 = now - 15 * 60 * 1000;
      const from30 = now - 30 * 60 * 1000;
      const from60 = now - 60 * 60 * 1000;

      if (f === '30m') {
        // setChartData(data.slice(0, -30 * 60));
        chart.xAxis[0].setExtremes(from30, now);
        setRangeValues({
          to: now,
          from: from30,
          range: 30 * 60 * 1000,
        });
      } else if (f === '1h') {
        // setChartData(data.slice(0, -60 * 60));
        chart.xAxis[0].setExtremes(from60, now);
        setRangeValues({
          to: now,
          from: from60,
          range: 60 * 60 * 1000,
        });
      } else {
        // setChartData(data.slice(0, -15 * 60));
        chart.xAxis[0].setExtremes(from15, now);
        setRangeValues({
          to: now,
          from: from15,
          range: 15 * 60 * 1000,
        });
      }
    }
  };

  return (
    <GraphTitleWrapper
      title={title}
      filter={filter}
      hasData={hasData}
      variant="subtitle1"
      handleFilterChange={(f: string) => {
        setFilter(f);
        // setRange(f);
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
