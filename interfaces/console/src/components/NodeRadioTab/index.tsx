/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import { RadioChartsConfig, TooltipsText } from '@/constants';
import { getMetricValue, isMetricValue } from '@/utils';
import { Grid, Paper, Stack } from '@mui/material';
import { useState } from 'react';
import LineChart from '../LineChart';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

const PLACEHOLDER_VALUE = 'NA';
interface INodeRadioTab {
  nodeId: string;
  metrics: any;
  loading: boolean;
  metricFrom: number;
}
const NodeRadioTab = ({
  nodeId,
  loading,
  metrics,
  metricFrom,
}: INodeRadioTab) => {
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
        <Paper
          sx={{
            p: 3,
            overflow: 'auto',
            height: { xs: 'calc(100vh - 480px)', md: 'calc(100vh - 328px)' },
          }}
        >
          <Stack spacing={4}>
            <LineChart
              nodeId={nodeId}
              loading={loading}
              topic={RadioChartsConfig[0].id}
              title={RadioChartsConfig[0].name}
              metricFrom={metricFrom}
              tabSection={Graphs_Type.Radio}
              hasData={isMetricValue(RadioChartsConfig[0].id, metrics)}
              initData={getMetricValue(RadioChartsConfig[0].id, metrics)}
            />
            <LineChart
              nodeId={nodeId}
              loading={loading}
              topic={RadioChartsConfig[1].id}
              title={RadioChartsConfig[1].name}
              metricFrom={metricFrom}
              tabSection={Graphs_Type.Radio}
              hasData={isMetricValue(RadioChartsConfig[1].id, metrics)}
              initData={getMetricValue(RadioChartsConfig[1].id, metrics)}
            />
            <LineChart
              nodeId={nodeId}
              loading={loading}
              topic={RadioChartsConfig[2].id}
              title={RadioChartsConfig[2].name}
              metricFrom={metricFrom}
              tabSection={Graphs_Type.Radio}
              hasData={isMetricValue(RadioChartsConfig[2].id, metrics)}
              initData={getMetricValue(RadioChartsConfig[2].id, metrics)}
            />
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeRadioTab;
