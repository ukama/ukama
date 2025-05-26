/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, NodeTypeEnum } from '@/client/graphql/generated';
import { MetricsStateRes } from '@/client/graphql/generated/subscriptions';
import { NODE_KPIS } from '@/constants';
import { getKPIStatValue, getMetricValue, isMetricValue } from '@/utils';
import { Paper, Stack } from '@mui/material';
import Grid from '@mui/material/Grid2';
import LineChart from '../LineChart';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

interface INodeResourcesTab {
  metrics: any;
  loading: boolean;
  metricFrom: number;
  statLoading: boolean;
  selectedNode: Node | undefined;
  nodeMetricsStatData: MetricsStateRes;
}
const NodeResourcesTab = ({
  metrics,
  loading,
  metricFrom,
  statLoading,
  nodeMetricsStatData,
}: INodeResourcesTab) => {
  const resourcesConfig = NODE_KPIS.RESOURCES[NodeTypeEnum.Tnode];

  return (
    <Grid container spacing={3}>
      <Grid size={{ xs: 12, md: 3 }}>
        <NodeStatsContainer
          index={0}
          selected={0}
          loading={statLoading}
          title={'Resources'}
          isCollapsable={false}
        >
          {resourcesConfig.map((config, i) => (
            <NodeStatItem
              id={config.id}
              name={config.name}
              unit={config.unit}
              format={config.format}
              key={`${config.id}-${i}`}
              threshold={config.threshold}
              nameInfo={config.description}
              value={getKPIStatValue(
                config.id,
                statLoading,
                nodeMetricsStatData,
              )}
            />
          ))}
        </NodeStatsContainer>
      </Grid>
      <Grid size={{ xs: 12, md: 9 }}>
        <Paper
          sx={{
            p: 3,
            overflow: 'auto',
            height: { xs: 'calc(100vh - 480px)', md: 'calc(100vh - 328px)' },
          }}
        >
          <Stack spacing={4}>
            {resourcesConfig.map((config, i) => (
              <LineChart
                from={metricFrom}
                topic={config.id}
                loading={loading}
                title={config.name}
                yunit={config.unit}
                format={config.format}
                key={`${config.id}-${i}`}
                tickInterval={config.tickInterval}
                tickPositions={config.tickPositions}
                hasData={isMetricValue(config.id, metrics)}
                initData={getMetricValue(config.id, metrics)}
              />
            ))}
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeResourcesTab;
