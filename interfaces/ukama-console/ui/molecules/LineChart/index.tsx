import { Box } from '@mui/material';
import { HighchartsReact } from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
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
      title: {
        style: {
          display: 'none',
        },
      },

      legend: { enabled: false },

      events: {
        load: function () {
          // var series: any = Highcharts.charts[0]?.series[0];
          // PubSub.subscribe(topic, (_, data) => {
          //   series.setData(data[0], true, true);
          // });
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
        data: (() => initData)(),
      },
    ],
    xAxis: {
      type: 'datetime',
      title: false,
      labels: {
        enabled: true,
        formatter: function (value: any) {
          return Highcharts.dateFormat('%H:%M:%S', value.value * 1000);
        },
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
