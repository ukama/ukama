/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Graphs_Type } from '@/client/graphql/generated/metrics';
import { TooltipsText } from '@/constants';
import { getMetricValue, isMetricValue } from '@/utils';
import { Grid, Paper, Stack } from '@mui/material';
import { useState } from 'react';
import LineChart from '../LineChart';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

const PLACEHOLDER_VALUE = 'NA';
interface INodeRadioTab {
  metrics: any;
  loading: boolean;
  metricFrom: number;
}
const NodeRadioTab = ({ loading, metrics, metricFrom }: INodeRadioTab) => {
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
              topic={'tx_power'}
              title={'TX Power'}
              metricFrom={metricFrom}
              tabSection={Graphs_Type.Radio}
              hasData={isMetricValue('tx_power', metrics)}
              initData={getMetricValue('tx_power', metrics)}
            />
            <LineChart
              loading={loading}
              topic={'rx_power'}
              title={'RX Power'}
              metricFrom={metricFrom}
              tabSection={Graphs_Type.Radio}
              hasData={isMetricValue('rx_power', metrics)}
              initData={getMetricValue('rx_power', metrics)}
            />
            <LineChart
              loading={loading}
              topic={'pa_power'}
              title={'PA Power'}
              metricFrom={metricFrom}
              tabSection={Graphs_Type.Radio}
              hasData={isMetricValue('pa_power', metrics)}
              initData={getMetricValue('pa_power', metrics)}
            />
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeRadioTab;
