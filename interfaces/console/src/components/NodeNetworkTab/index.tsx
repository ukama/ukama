/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { MetricsRes } from '@/client/graphql/generated/subscriptions';
import { NetworkChartsConfig, TooltipsText } from '@/constants';
import { getMetricValue, isMetricValue } from '@/utils';
import { Grid, Paper, Stack } from '@mui/material';
import { useState } from 'react';
import LineChart from '../LineChart';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

const PLACEHOLDER_VALUE = 'NA';
interface INodeOverviewTab {
  nodeId: string;
  metrics: MetricsRes;
  metricFrom: number;
  loading: boolean;
}
const NodeNetworkTab = ({
  nodeId,
  loading,
  metrics,
  metricFrom,
}: INodeOverviewTab) => {
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
            name={NetworkChartsConfig[0].name}
            nameInfo={TooltipsText.DL}
          />
          <NodeStatItem
            variant={'large'}
            value={PLACEHOLDER_VALUE}
            name={NetworkChartsConfig[1].name}
            nameInfo={TooltipsText.UL}
          />
          <NodeStatItem
            value={PLACEHOLDER_VALUE}
            variant={'large'}
            name={NetworkChartsConfig[2].name}
            nameInfo={TooltipsText.RRCCNX}
          />
          <NodeStatItem
            value={PLACEHOLDER_VALUE}
            variant={'large'}
            name={NetworkChartsConfig[3].name}
            nameInfo={TooltipsText.ERAB}
          />
          <NodeStatItem
            value={PLACEHOLDER_VALUE}
            variant={'large'}
            name={NetworkChartsConfig[4].name}
            nameInfo={TooltipsText.RLS}
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
              loading={loading}
              from={metricFrom}
              topic={NetworkChartsConfig[0].id}
              title={NetworkChartsConfig[0].name}
              initData={getMetricValue(NetworkChartsConfig[0].id, metrics)}
              hasData={isMetricValue(NetworkChartsConfig[0].id, metrics)}
            />
            <LineChart
              loading={loading}
              from={metricFrom}
              title={NetworkChartsConfig[1].name}
              topic={NetworkChartsConfig[1].id}
              initData={getMetricValue(NetworkChartsConfig[1].id, metrics)}
              hasData={isMetricValue(NetworkChartsConfig[1].id, metrics)}
            />
            <LineChart
              loading={loading}
              from={metricFrom}
              topic={NetworkChartsConfig[2].id}
              title={NetworkChartsConfig[2].name}
              initData={getMetricValue(NetworkChartsConfig[2].id, metrics)}
              hasData={isMetricValue(NetworkChartsConfig[2].id, metrics)}
            />
            <LineChart
              loading={loading}
              from={metricFrom}
              topic={NetworkChartsConfig[3].id}
              title={NetworkChartsConfig[3].name}
              initData={getMetricValue(NetworkChartsConfig[3].id, metrics)}
              hasData={isMetricValue(NetworkChartsConfig[3].id, metrics)}
            />
            <LineChart
              loading={loading}
              from={metricFrom}
              topic={NetworkChartsConfig[4].id}
              title={NetworkChartsConfig[4].name}
              initData={getMetricValue(NetworkChartsConfig[4].id, metrics)}
              hasData={isMetricValue(NetworkChartsConfig[4].id, metrics)}
            />
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeNetworkTab;
