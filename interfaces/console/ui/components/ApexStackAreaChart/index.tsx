import { format } from 'date-fns';
import dynamic from 'next/dynamic';
import React from 'react';
import { GraphTitleWrapper } from '..';
const Chart = dynamic(() => import('react-apexcharts'), { ssr: false });

const TIME_RANGE_IN_MILLISECONDS = 100;

interface IApexLineChartIntegration {
  data: any;
  name: string;
  filter?: string;
  hasData: boolean;
  refreshInterval?: number;
  onRefreshData?: Function;
  onFilterChange?: Function;
}

const ApexLineChart = (props: any) => {
  const options: any = {
    series: [
      {
        name: 'series1',
        data: [31, 40, 28, 51, 42, 109, 100],
      },
    ],
    options: {
      chart: {
        height: 350,
        type: 'area',
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        curve: 'smooth',
      },
      xaxis: {
        type: 'datetime',
        range: props.range,
        labels: {
          formatter: (val: any) =>
            val ? format(new Date(val * 1000), 'mm:ss') : '',
        },
        tooltip: {
          enabled: false,
          offsetX: 0,
        },
      },
      yaxis: {
        labels: {
          formatter: (val: any) => val.toFixed(2),
        },
        // min: 0,
        // max: 100,
        // tooltip: {
        //     enabled: true,
        // },
        tickAmount: 8,
      },
    },
  };

  return (
    <Chart
      type="area"
      key={props.name}
      height={'300px'}
      options={options}
      series={props.dataList}
    />
  );
};

const ApexStackChartIntegration = ({
  name,
  data = [],
  onRefreshData,
  filter = 'LIVE',
  hasData = false,
  refreshInterval = 10000,
  onFilterChange = () => {
    /*DEFAULT FUNCTION*/
  },
}: IApexLineChartIntegration) => {
  React.useEffect(() => {
    const interval = setInterval(() => {
      onRefreshData && onRefreshData();
    }, refreshInterval);

    return () => clearInterval(interval);
  });

  return (
    <GraphTitleWrapper
      key={name}
      title={name}
      filter={filter}
      hasData={hasData}
      variant="subtitle1"
      handleFilterChange={onFilterChange}
    >
      <ApexLineChart
        key={name}
        name={name}
        dataList={data}
        range={TIME_RANGE_IN_MILLISECONDS}
      />
    </GraphTitleWrapper>
  );
};

export default ApexStackChartIntegration;
