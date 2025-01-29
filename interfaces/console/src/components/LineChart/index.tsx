/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import { useAppContext } from '@/context';
import { Box } from '@mui/material';
import { HighchartsReact } from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import PubSub from 'pubsub-js';
import GraphTitleWrapper from '../GraphTitleWrapper';
import { useEffect } from 'react';

interface ILineChart {
  nodeId: string;
  metricFrom: any;
  topic: string;
  initData: any;
  title?: string;
  filter?: string;
  hasData?: boolean;
  loading?: boolean;
  tabSection: Graphs_Type;
}

const getOptions = (topic: string, title: string, initData: any) => {
  return {
    title: {
      text: topic,
      align: 'left',
    },
    chart: {
      legend: { enabled: false },
      events: {
        load: function () {
          const chart: any =
            Highcharts.charts.length > 0
              ? Highcharts.charts.find((c: any) => c?.title?.textStr === topic)
              : null;
          if (chart) {
            const series: any = chart?.series[0];
            PubSub.subscribe(topic, (_, data) => {
              if (topic === chart?.title?.textStr && series) {
                series.addPoint(data, true, true);
              }
            });
          }
        },
      },
    },

    time: {
      useUTC: false,
    },

    rangeSelector: {
      buttons: [
        {
          count: 30,
          type: 'second',
          text: '30S',
        },
        {
          count: 1,
          type: 'minute',
          text: '1M',
        },
        {
          type: 'all',
          text: 'All',
        },
      ],
      inputEnabled: false,
      selected: 0,
    },

    exporting: {
      enabled: true,
    },

    navigator: {
      enabled: false,
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
        formate: '{value:%H:%M:%S}',
      },
    },
    yAxis: {
      title: false,
      opposite: false,
    },
  };
};

const LineChart = ({
  topic,
  nodeId,
  hasData,
  initData,
  metricFrom,
  title = '',
  loading = false,
  filter = 'LIVE',
  tabSection = Graphs_Type.NodeHealth,
}: ILineChart) => {
  useEffect(() => {
    console.log('LineChart mounting for:', {
      topic,
      nodeId,
      hasData,
      dataLength: initData?.length,
      initData: initData?.slice(0, 2), // Log first two items for debugging
    });
  }, []);

  const defaultData = [
    [Date.now() - 60000, 0],
    [Date.now(), 0],
  ];

  const chartOptions = getOptions(
    topic,
    title,
    initData?.length > 0 ? initData : defaultData,
  );

  return (
    <GraphTitleWrapper
      filter={filter}
      hasData={hasData}
      variant="subtitle1"
      title={title}
      handleFilterChange={() => {}}
      loading={loading}
    >
      <Box sx={{ width: '100%' }}>
        <HighchartsReact
          key={`${topic}-${nodeId}-${initData?.length}`}
          options={chartOptions}
          highcharts={Highcharts}
        />
      </Box>
    </GraphTitleWrapper>
  );
};

export default LineChart;
