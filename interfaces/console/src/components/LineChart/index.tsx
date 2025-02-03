/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import { Box } from '@mui/material';
import { HighchartsReact } from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import PubSub from 'pubsub-js';
import { useEffect } from 'react';
import GraphTitleWrapper from '../GraphTitleWrapper';

interface ILineChart {
  nodeId?: string;
  siteId?: string;
  metricFrom: any;
  topic: string;
  initData: any;
  title?: string;
  filter?: string;
  hasData?: boolean;
  loading?: boolean;
  tabSection: Graphs_Type;
}

const getOptions = (
  topic: string,
  title: string,
  initData: any,
  entityId: string,
) => {
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
              ? Highcharts.charts.find(
                  (c: any) =>
                    c?.title?.textStr === topic &&
                    c?.userOptions?.customData?.entityId === entityId,
                )
              : null;
          if (chart) {
            const series: any = chart?.series[0];
            PubSub.subscribe(topic, (_, data) => {
              if (
                topic === chart?.title?.textStr &&
                series &&
                ((data.nodeId && data.nodeId === entityId) ||
                  (data.siteId && data.siteId === entityId))
              ) {
                series.addPoint(data, true, true);
              }
            });
          }
        },
      },
    },
    customData: {
      entityId,
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
  siteId,
  hasData,
  initData,
  metricFrom,
  title = '',
  loading = false,
  filter = 'LIVE',
  tabSection = Graphs_Type.NodeHealth,
}: ILineChart) => {
  useEffect(() => {
    const entityType = siteId ? 'site' : 'node';
    const entityId = siteId || nodeId;

    console.log('LineChart mounting for:', {
      topic,
      entityType,
      entityId,
      hasData,
      dataLength: initData?.length,
      initData: initData?.slice(0, 2),
    });

    return () => {
      PubSub.unsubscribe(topic);
    };
  }, [topic, nodeId, siteId, initData]);

  const defaultData = [
    [Date.now() - 60000, 0],
    [Date.now(), 0],
  ];

  const entityId = siteId || nodeId;

  const chartOptions = getOptions(
    topic,
    title,
    initData?.length > 0 ? initData : defaultData,
    entityId as string,
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
      <Box
        sx={{
          width: '100%',
        }}
      >
        <HighchartsReact
          key={`${topic}-${entityId}-${initData?.length}`}
          options={{
            ...chartOptions,
            chart: {
              ...chartOptions.chart,
            },
          }}
          highcharts={Highcharts}
        />
      </Box>
    </GraphTitleWrapper>
  );
};

export default LineChart;
