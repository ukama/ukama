import { Graphs_Type } from '@/generated/metrics';
import { Box } from '@mui/material';
import { HighchartsReact } from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import PubSub from 'pubsub-js';
import GraphTitleWrapper from '../GraphTitleWrapper';
import MetricSubscription from './metricSubscription';

interface ILineChart {
  metricFrom: any;
  topic: string;
  initData: any;
  title?: string;
  filter?: string;
  hasData?: boolean;
  loading?: boolean;
  tabSection: Graphs_Type;
  onFilterChange?: Function;
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
          var chart: any =
            Highcharts.charts.length > 0
              ? Highcharts.charts.find((c: any) => c?.title?.textStr === topic)
              : null;
          if (chart) {
            var series: any = chart?.series[0];
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
};

const LineChart = ({
  topic,
  hasData,
  initData,
  metricFrom,
  title = '',
  loading = false,
  filter = 'LIVE',
  tabSection = Graphs_Type.NodeHealth,
}: ILineChart) => {
  return (
    <GraphTitleWrapper
      filter={filter}
      hasData={hasData}
      variant="subtitle1"
      title={title}
      handleFilterChange={() => {}}
      loading={loading || !initData}
    >
      <Box sx={{ width: '100%' }}>
        <MetricSubscription type={tabSection} from={metricFrom} />
        <HighchartsReact
          key={topic}
          options={getOptions(topic, title, initData)}
          highcharts={Highcharts}
        />
      </Box>
    </GraphTitleWrapper>
  );
};

export default LineChart;
