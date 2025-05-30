/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, NodeTypeEnum } from '@/client/graphql/generated';
import {
  Graphs_Type,
  MetricsRes,
  MetricsStateRes,
} from '@/client/graphql/generated/subscriptions';
import { NODE_KPIS } from '@/constants';
import { getKPIStatValue, getMetricValue, isMetricValue } from '@/utils';
import { Paper, Stack } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { useState } from 'react';
import LineChart from '../LineChart';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

interface INodeOverviewTab {
  loading: boolean;
  metricFrom: number;
  metrics: MetricsRes;
  statLoading: boolean;
  selectedNode: Node | undefined;
  handleSectionChange: (section: Graphs_Type) => void;
  nodeMetricsStatData: MetricsStateRes;
}
const NodeNetworkTab = ({
  loading,
  metrics,
  metricFrom,
  statLoading,
  handleSectionChange,
  nodeMetricsStatData,
}: INodeOverviewTab) => {
  const [selected, setSelected] = useState<number>(0);
  const networkCellular = NODE_KPIS.NETWORK_CELLULAR[NodeTypeEnum.Tnode];
  const networkBackhaul = NODE_KPIS.NETWORK_BACKHAUL[NodeTypeEnum.Tnode];

  const handleOnSelected = (value: number) => {
    handleSectionChange(
      value === 1 ? Graphs_Type.NetworkBackhaul : Graphs_Type.NetworkCellular,
    );
    setSelected((prev) => (prev === 1 ? 0 : 1));
  };
  return (
    <Grid container spacing={3}>
      <Grid size={{ xs: 12, md: 3 }}>
        <Stack spacing={2}>
          <NodeStatsContainer
            index={0}
            title={'Cellular'}
            isClickable={true}
            selected={selected}
            loading={statLoading}
            isCollapsable={false}
            handleAction={handleOnSelected}
          >
            {networkCellular.map((config, i) => (
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
          <NodeStatsContainer
            index={1}
            loading={statLoading}
            isClickable={true}
            selected={selected}
            title={'Backhaul'}
            handleAction={handleOnSelected}
          >
            {networkBackhaul.map((config, i) => (
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
        </Stack>
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
            {selected === 1
              ? networkBackhaul.map((config, i) => (
                  <LineChart
                    from={metricFrom}
                    topic={config.id}
                    loading={loading}
                    yunit={config.unit}
                    title={config.name}
                    format={config.format}
                    key={`${config.id}-${i}`}
                    tickInterval={config.tickInterval}
                    tickPositions={config.tickPositions}
                    hasData={isMetricValue(config.id, metrics)}
                    initData={getMetricValue(config.id, metrics)}
                  />
                ))
              : networkCellular.map((config, i) => (
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

export default NodeNetworkTab;
