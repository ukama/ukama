import { Box } from '@mui/material';
import { HighchartsReact } from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import PubSub from 'pubsub-js';
import GraphTitleWrapper from '../GraphTitleWrapper';
import { MetricSub } from './metricSub';

interface ILineChart {
  metricFrom: any;
  topic: string;
  initData: any;
  title?: string;
  filter?: string;
  hasData?: boolean;
  loading?: boolean;
  onFilterChange?: Function;
}

const LineChart = ({
  title,
  topic,
  initData,
  metricFrom,
  loading = false,
  filter = 'LIVE',
}: ILineChart) => {
  const options = {
    chart: {
      title: {
        style: {
          display: 'none',
        },
      },

      legend: { enabled: false },

      events: {
        load: function () {
          var series: any = Highcharts.charts[0]?.series[0];
          PubSub.subscribe('memory_trx_used', (_, data) => {
            console.log('Subscription data: ', data);
            series.addPoint(data, true, true);
          });
          // setInterval(function () {
          //   var x = new Date().getTime() / 1000, // current time
          //     y = Math.round(Math.random() * 100);
          //   series.addPoint([Math.floor(x) * 1000, y], true, true);
          // }, 1000);
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

    series: [
      {
        name: title,
        data: (function () {
          var data = [...initData];
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

  return (
    <GraphTitleWrapper
      filter={filter}
      variant="subtitle1"
      title={title || ''}
      handleFilterChange={() => {}}
      loading={loading || !initData}
      hasData={initData?.length > 0 || false}
    >
      <Box sx={{ width: '100%' }}>
        <MetricSub from={metricFrom} />
        <HighchartsReact options={options} highcharts={Highcharts} />
      </Box>
    </GraphTitleWrapper>
  );
};

export default LineChart;
