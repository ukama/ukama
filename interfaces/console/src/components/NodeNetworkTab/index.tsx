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

const PLACEHOLDER_VALUE = 'NA';
interface INodeOverviewTab {
  loading: boolean;
  metricFrom: number;
  metrics: MetricsRes;
  selectedNode: Node | undefined;
  handleSectionChange: Function;
  nodeMetricsStatData: MetricsStateRes;
}
const NodeNetworkTab = ({
  loading,
  metrics,
  metricFrom,
  selectedNode,
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
    setSelected(value);
  };

  return (
    <Grid container spacing={3}>
      <Grid size={{ xs: 12, md: 3 }}>
        <Stack spacing={2}>
          <NodeStatsContainer
            index={0}
            loading={loading}
            title={'Cellular'}
            selected={selected}
            isCollapsable={false}
            handleAction={handleOnSelected}
          >
            {networkCellular.map((config, i) => (
              <NodeStatItem
                id={config.id}
                name={config.name}
                unit={config.unit}
                key={`${config.id}-${i}`}
                nameInfo={config.description}
                value={getKPIStatValue(config.id, loading, nodeMetricsStatData)}
              />
            ))}
          </NodeStatsContainer>
          <NodeStatsContainer
            index={1}
            loading={loading}
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
                key={`${config.id}-${i}`}
                nameInfo={config.description}
                value={getKPIStatValue(config.id, loading, nodeMetricsStatData)}
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
                    title={config.name}
                    key={`${config.id}-${i}`}
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
                    key={`${config.id}-${i}`}
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
