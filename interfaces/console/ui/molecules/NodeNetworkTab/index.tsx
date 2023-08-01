import { TooltipsText } from '@/constants';
import { Grid, Paper, Stack } from '@mui/material';
import { useState } from 'react';
import LineChart from '../LineChart';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

const PLACEHOLDER_VALUE = 'NA';
interface INodeOverviewTab {
  metrics: any;
  loading: boolean;
}
const NodeNetworkTab = ({ loading, metrics }: INodeOverviewTab) => {
  const [isCollapse, setIsCollapse] = useState<boolean>(false);
  const handleCollapse = () => setIsCollapse((prev) => !prev);

  return (
    <Grid container spacing={3}>
      <Grid md xs item lg={!isCollapse ? 4 : 1}>
        <NodeStatsContainer
          index={0}
          selected={0}
          loading={loading}
          title={'Network'}
          isCollapsable={true}
          isCollapse={isCollapse}
          onCollapse={handleCollapse}
        >
          <NodeStatItem
            variant={'large'}
            value={PLACEHOLDER_VALUE}
            name={'Throughput (D/L)'}
            nameInfo={TooltipsText.DL}
          />
          <NodeStatItem
            variant={'large'}
            value={PLACEHOLDER_VALUE}
            name={'Throughput (U/L)'}
            nameInfo={TooltipsText.UL}
          />
          <NodeStatItem
            value={PLACEHOLDER_VALUE}
            variant={'large'}
            name={'RRC CNX Success'}
            nameInfo={TooltipsText.RRCCNX}
          />
          <NodeStatItem
            value={PLACEHOLDER_VALUE}
            variant={'large'}
            name={'ERAB Drop Rate'}
            nameInfo={TooltipsText.ERAB}
          />
          <NodeStatItem
            value={PLACEHOLDER_VALUE}
            variant={'large'}
            name={'RLS  Drop Rate'}
            nameInfo={TooltipsText.RLS}
          />
        </NodeStatsContainer>
      </Grid>
      <Grid item lg={isCollapse ? 11 : 8} md xs>
        <Paper sx={{ p: 3, width: '100%' }}>
          <Stack spacing={4}>
            <LineChart
              loading={loading}
              initData={metrics}
              hasData={metrics.length > 0}
              topic={'throughputuplink'}
              title={'Throughput (U/L)'}
            />
            <LineChart
              loading={loading}
              initData={metrics}
              hasData={metrics.length > 0}
              topic={'throughputdownlink'}
              title={'Throughput (D/L)'}
            />
            <LineChart
              loading={loading}
              initData={metrics}
              hasData={metrics.length > 0}
              topic={'rrc'}
              title={'RRC'}
            />
            <LineChart
              loading={loading}
              initData={metrics}
              hasData={metrics.length > 0}
              topic={'erab'}
              title={'ERAB'}
            />
            <LineChart
              loading={loading}
              initData={metrics}
              hasData={metrics.length > 0}
              topic={'rlc'}
              title={'RLC'}
            />
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeNetworkTab;
