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
import {
  getKPIStatValue,
  getMetricValue,
  isMetricValue,
  NodeEnumToString,
} from '@/utils';
import { Paper, Stack } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { useEffect, useState } from 'react';
import LineChart from '../LineChart';
import NodeDetailsCard from '../NodeDetailsCard';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

interface INodeOverviewTab {
  nodeId: string;
  metricFrom: number;
  metrics: MetricsRes;
  statLoading: boolean;
  metricsLoading: boolean;
  onNodeSelected: (node: Node) => void;
  isUpdateAvailable: boolean;
  handleUpdateNode: () => void;
  selectedNode: Node | undefined;
  connectedUsers: string | undefined;
  getNodeSoftwareUpdateInfos: () => void;
  handleOverviewSectionChange: (section: Graphs_Type) => void;
  nodeMetricsStatData: MetricsStateRes;
}

const NodeOverviewTab = ({
  metrics,
  metricFrom,
  statLoading,
  selectedNode,
  metricsLoading,
  handleUpdateNode,
  isUpdateAvailable,
  nodeMetricsStatData,
  getNodeSoftwareUpdateInfos,
  handleOverviewSectionChange,
}: INodeOverviewTab) => {
  const nodeType = selectedNode?.type ?? NodeTypeEnum.Tnode;
  const healthConfig = NODE_KPIS.HEALTH[NodeTypeEnum.Tnode];
  const subscriberConfig = NODE_KPIS.SUBSCRIBER[NodeTypeEnum.Tnode];
  const [selected, setSelected] = useState<number>(0);
  useEffect(() => {
    setSelected(0);
  }, [selectedNode]);

  const handleOnSelected = (value: number) => {
    handleOverviewSectionChange(
      value === 1 ? Graphs_Type.NodeHealth : Graphs_Type.Subscribers,
    );
    setSelected(value);
  };

  return (
    <Grid container columnSpacing={3} rowSpacing={2}>
      <Grid size={{ xs: 12, md: 3.5 }}>
        <Stack spacing={2}>
          <NodeStatsContainer
            index={0}
            isClickable={true}
            selected={selected}
            loading={statLoading}
            title={'Node Information'}
            handleAction={handleOnSelected}
          >
            <NodeStatItem
              name={'Model type'}
              value={NodeEnumToString(nodeType)}
            />

            <NodeStatItem
              name={'Serial #'}
              value={selectedNode?.id.toLowerCase() ?? '-'}
            />
            {/* {selectedNode?.type === 'TOWER' && (
                <Grid item xs={12}>
                  <NodeGroup
                    nodes={nodeGroupData?.attached ?? []}
                    loading={nodeGroupLoading}
                    handleNodeAction={onNodeSelected}
                  />
                </Grid>
              )} */}
          </NodeStatsContainer>
          <NodeStatsContainer
            index={1}
            title={'Health'}
            isClickable={true}
            selected={selected}
            loading={statLoading}
            handleAction={handleOnSelected}
          >
            {healthConfig.map((config, i) => (
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
          {selectedNode?.type !== NodeTypeEnum.Anode && (
            <NodeStatsContainer
              index={2}
              loading={statLoading}
              isClickable={true}
              selected={selected}
              title={'Subscribers'}
              handleAction={handleOnSelected}
            >
              {subscriberConfig.map((config, i) => (
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
          )}
        </Stack>
      </Grid>
      <Grid size={{ xs: 12, md: 8.5 }}>
        {selected === 0 && (
          <NodeDetailsCard
            nodeType={selectedNode?.type ?? undefined}
            getNodeUpdateInfos={getNodeSoftwareUpdateInfos}
            loading={statLoading}
            nodeTitle={selectedNode?.name ?? 'HOME'}
            handleUpdateNode={handleUpdateNode}
            isUpdateAvailable={isUpdateAvailable}
          />
        )}
        <Paper
          sx={{
            p: 3,
            overflow: 'auto',
            display: selected === 0 ? 'none' : 'block',
            height: { xs: 'calc(100vh - 480px)', md: 'calc(100vh - 328px)' },
          }}
        >
          {selected === 1 && (
            <Stack spacing={4}>
              {healthConfig.map((config, i) => (
                <LineChart
                  from={metricFrom}
                  topic={config.id}
                  title={config.name}
                  yunit={config.unit}
                  format={config.format}
                  loading={metricsLoading}
                  key={`${config.id}-${i}`}
                  tickInterval={config.tickInterval}
                  tickPositions={config.tickPositions}
                  hasData={isMetricValue(config.id, metrics)}
                  initData={getMetricValue(config.id, metrics)}
                />
              ))}
            </Stack>
          )}
          {selected === 2 && nodeType === NodeTypeEnum.Tnode && (
            <Stack spacing={4}>
              {subscriberConfig.map((config, i) => (
                <LineChart
                  from={metricFrom}
                  topic={config.id}
                  title={config.name}
                  yunit={config.unit}
                  format={config.format}
                  loading={metricsLoading}
                  key={`${config.id}-${i}`}
                  tickInterval={config.tickInterval}
                  tickPositions={config.tickPositions}
                  hasData={isMetricValue(config.id, metrics)}
                  initData={getMetricValue(config.id, metrics)}
                />
              ))}
            </Stack>
          )}
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeOverviewTab;
