import { Box } from '@mui/material';
import { HighchartsReact } from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import PubSub from 'pubsub-js';
import GraphTitleWrapper from '../GraphTitleWrapper';

interface ILineChart {
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
  loading = false,
  filter = 'LIVE',
}: ILineChart) => {
  const options = {
    chart: {
      events: {
        load: function () {
          var series: any = Highcharts.charts[0]?.series[0];
          PubSub.subscribe(topic, (_, data) => {
            series.setData(data[0], true, true);
          });
          //  setInterval(function () {
          //   var x = new Date().getTime(),
          //    y = Math.round(Math.random() * 100)
          //   series.addPoint([x, y], true, true)
          //  }, 1000)
        },
      },
    },

    time: {
      useUTC: false,
    },

    rangeSelector: {
      buttons: [
        {
          count: 1,
          type: 'minute',
          text: '1M',
        },
        {
          count: 5,
          type: 'minute',
          text: '5M',
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
        data: (() => initData)(),
      },
    ],
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
        <HighchartsReact
          options={options}
          highcharts={Highcharts}
          constructorType={'stockChart'}
        />
      </Box>
    </GraphTitleWrapper>
  );
};

export default LineChart;
