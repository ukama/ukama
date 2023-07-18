import { isDarkmode } from '@/app-recoil';
import { colors } from '@/styles/theme';
import { format } from 'date-fns';
import dynamic from 'next/dynamic';
import { useRecoilValue } from 'recoil';
const Chart = dynamic(() => import('react-apexcharts'), { ssr: false });

interface ILineChart {
  range: any;
  name: string;
  dataList: any;
}

const LineChart = ({ name, dataList, range }: ILineChart) => {
  const _isDarkMod = useRecoilValue(isDarkmode);
  const options: any = {
    stroke: {
      lineCap: 'butt',
      curve: 'smooth',
      width: 5,
    },
    chart: {
      minHeight: '200px',
      height: '100%',
      width: '100%',
      zoom: {
        type: 'x',
        enabled: false,
        autoScaleYaxis: true,
      },
      animations: {
        enabled: true,
        easing: 'linear',
        dynamicAnimation: {
          speed: 1000,
        },
      },
      dropShadow: {
        enabled: true,
        top: 1,
        left: 1,
        bottom: 1,
        blur: 3,
        opacity: 0.2,
      },
      toolbar: {
        show: false,
        tools: {
          download: false,
        },
      },
    },
    grid: {
      borderColor: _isDarkMod ? colors.vulcan60 : colors.white60,
      opacity: 0.3,
    },
    tooltip: {
      theme: _isDarkMod ? 'dark' : 'light',
      y: {
        formatter: (val: any) => val.toFixed(8),
      },
    },
    xaxis: {
      type: 'datetime',
      range: range,
      labels: {
        formatter: (val: any) =>
          val ? format(new Date(val * 1000), 'hh:mm:ss') : '',
      },
      tooltip: {
        enabled: false,
        offsetX: 0,
      },
    },

    yaxis: {
      labels: {
        formatter: (val: any) => val.toFixed(4),
      },
      // min: 0,
      // max: 100,
      // tooltip: {
      //     enabled: true,
      // },
      tickAmount: 8,
    },
  };
  return (
    <Chart
      type="line"
      key={name}
      options={options}
      height={'300px'}
      series={dataList}
    />
  );
};

export default LineChart;
