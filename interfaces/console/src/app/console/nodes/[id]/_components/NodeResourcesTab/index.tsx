/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, NodeTypeEnum } from '@/client/graphql/generated';
import { MetricsRes, MetricsStateRes } from '@/client/graphql/generated/subscriptions';
import { NODE_KPIS } from '@/constants';
import { KpiConfig } from '@/types';
import { getKPIStatValue, getMetricValue, isMetricValue } from '@/utils';
import { Paper, Stack } from '@mui/material';
import Grid from '@mui/material/Grid2';
import LineChart from '@/components/ui/LineChart';
import NodeStatItem from '@/app/console/nodes/[id]/_components/NodeStatItem';
import NodeStatsContainer from '@/app/console/nodes/[id]/_components/NodeStatsContainer';

const withApiDisplayMeta = (
  config: KpiConfig,
  statsData: MetricsStateRes,
  metrics: MetricsRes,
) => {
  const statMeta = statsData?.metrics?.find((m) => m.type === config.id);
  const rangeMeta = metrics?.metrics?.find((m) => m.type === config.id);
  const meta = statMeta ?? rangeMeta;

  return {
    ...config,
    unit: meta?.unit ?? config.unit,
    format: meta?.format ?? config.format ?? 'number',
    tickInterval: meta?.tickInterval ?? config.tickInterval,
    tickPositions: meta?.tickPositions ?? config.tickPositions,
    threshold: meta?.threshold ?? config.threshold,
  };
};

interface INodeResourcesTab {
  metrics: MetricsRes;
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
  const resourcesConfig = NODE_KPIS.RESOURCES[NodeTypeEnum.Tnode].map((config) =>
    withApiDisplayMeta(config, nodeMetricsStatData, metrics),
  );

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
