import LineChart from './LineChart';
import { GraphTitleWrapper } from '..';
import { makeStyles } from '@mui/styles';
import { Box, Theme } from '@mui/material';

const TIME_RANGE_IN_MILLISECONDS = 100;

const useStyles = makeStyles<Theme>((theme) => ({
  chartStyle: {
    width: '100%',
    height: '100%',
    '& .apexcharts-yaxis-label tspan': {
      fill: theme.palette.text.primary,
    },
    '& .apexcharts-xaxis-label tspan': {
      fill: theme.palette.text.primary,
    },
  },
}));

interface IApexLineChartIntegration {
  data: any;
  name?: string;
  filter?: string;
  hasData?: boolean;
  loading?: boolean;
  onFilterChange?: Function;
}

const ApexLineChart = ({
  data,
  loading = false,
  filter = 'LIVE',
  onFilterChange = () => {
    /*DEFAULT FUNCTION*/
  },
}: IApexLineChartIntegration) => {
  const classes = useStyles();
  return (
    <GraphTitleWrapper
      filter={filter}
      loading={loading || !data}
      variant="subtitle1"
      title={data?.name || ''}
      handleFilterChange={onFilterChange}
      hasData={data?.data.length > 0 || false}
    >
      <Box component="div" className={classes.chartStyle}>
        <LineChart
          name={data?.name || ''}
          dataList={[data]}
          range={TIME_RANGE_IN_MILLISECONDS}
        />
      </Box>
    </GraphTitleWrapper>
  );
};

export default ApexLineChart;
