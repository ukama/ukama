/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Graphs_Type, MetricsRes } from '@/client/graphql/generated/metrics';
import { TooltipsText } from '@/constants';
import { getMetricValue, isMetricValue } from '@/utils';
import { Grid, Paper, Stack } from '@mui/material';
import { useState } from 'react';
import LineChart from '../LineChart';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

const PLACEHOLDER_VALUE = 'NA';
interface INodeOverviewTab {
  metrics: MetricsRes;
  metricFrom: number;
  loading: boolean;
}
const NodeNetworkTab = ({ loading, metrics, metricFrom }: INodeOverviewTab) => {
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
              metricFrom={metricFrom}
              topic={'throughputuplink'}
              title={'Throughput (U/L)'}
              tabSection={Graphs_Type.Network}
              initData={getMetricValue('throughputuplink', metrics)}
              hasData={isMetricValue('throughputuplink', metrics)}
            />
            <LineChart
              loading={loading}
              metricFrom={metricFrom}
              topic={'throughputdownlink'}
              title={'Throughput (D/L)'}
              tabSection={Graphs_Type.Network}
              initData={getMetricValue('throughputdownlink', metrics)}
              hasData={isMetricValue('throughputdownlink', metrics)}
            />
            <LineChart
              topic={'rrc'}
              title={'RRC'}
              loading={loading}
              metricFrom={metricFrom}
              tabSection={Graphs_Type.Network}
              initData={getMetricValue('rrc', metrics)}
              hasData={isMetricValue('rrc', metrics)}
            />
            <LineChart
              topic={'erab'}
              title={'ERAB'}
              loading={loading}
              metricFrom={metricFrom}
              tabSection={Graphs_Type.Network}
              initData={getMetricValue('erab', metrics)}
              hasData={isMetricValue('erab', metrics)}
            />
            <LineChart
              topic={'rlc'}
              title={'RLC'}
              loading={loading}
              metricFrom={metricFrom}
              tabSection={Graphs_Type.Network}
              initData={getMetricValue('rlc', metrics)}
              hasData={isMetricValue('rlc', metrics)}
            />
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeNetworkTab;
