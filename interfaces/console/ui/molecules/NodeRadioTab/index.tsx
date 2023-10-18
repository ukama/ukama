import { TooltipsText } from '@/constants';
import { Grid, Paper, Stack } from '@mui/material';
import { useState } from 'react';
import LineChart from '../LineChart';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

const PLACEHOLDER_VALUE = 'NA';
interface INodeRadioTab {
  metrics: any;
  loading: boolean;
}
const NodeRadioTab = ({ loading, metrics }: INodeRadioTab) => {
  const [isCollapse, setIsCollapse] = useState<boolean>(false);
  const handleCollapse = () => setIsCollapse((prev) => !prev);
  return (
    <Grid container spacing={3}>
      <Grid item lg={!isCollapse ? 4 : 1} md xs>
        <NodeStatsContainer
          index={0}
          selected={0}
          title={'Radio'}
          loading={loading}
          isCollapsable={true}
          isCollapse={isCollapse}
          onCollapse={handleCollapse}
        >
          <NodeStatItem
            value={PLACEHOLDER_VALUE}
            variant={'large'}
            name={'TX Power'}
            nameInfo={TooltipsText.TXPOWER}
          />
          <NodeStatItem
            value={PLACEHOLDER_VALUE}
            variant={'large'}
            name={'RX Power'}
            nameInfo={TooltipsText.RXPOWER}
          />
          <NodeStatItem
            value={PLACEHOLDER_VALUE}
            name={'PA Power'}
            variant={'large'}
            nameInfo={TooltipsText.PAPOWER}
          />
        </NodeStatsContainer>
      </Grid>
      <Grid item lg={isCollapse ? 11 : 8} md xs>
        <Paper sx={{ p: 3, width: '100%' }}>
          <Stack spacing={4}>
            <LineChart
              loading={loading}
              initData={metrics}
              metricFrom={0}
              hasData={metrics.length > 0}
              topic={'txpower'}
              title={'TX Power'}
            />
            <LineChart
              loading={loading}
              initData={metrics}
              metricFrom={0}
              hasData={metrics.length > 0}
              topic={'rxpower'}
              title={'RX Power'}
            />
            <LineChart
              loading={loading}
              initData={metrics}
              metricFrom={0}
              hasData={metrics.length > 0}
              topic={'papower'}
              title={'PA Power'}
            />
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeRadioTab;
