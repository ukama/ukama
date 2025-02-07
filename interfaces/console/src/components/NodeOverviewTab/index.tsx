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
} from '@/client/graphql/generated/subscriptions';
import { KPI_PLACEHOLDER_VALUE, NODE_KPIS } from '@/constants';
import { getMetricValue, isMetricValue } from '@/utils';
import { Paper, Stack, capitalize } from '@mui/material';
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
  loading: boolean;
  metricsLoading: boolean;
  onNodeSelected: Function;
  isUpdateAvailable: boolean;
  handleUpdateNode: Function;
  selectedNode: Node | undefined;
  connectedUsers: string | undefined;
  getNodeSoftwareUpdateInfos: Function;
  handleOverviewSectionChange: Function;
}

const NodeOverviewTab = ({
  nodeId,
  metrics,
  loading,
  metricFrom,
  selectedNode,
  metricsLoading,
  connectedUsers = '0',
  onNodeSelected,
  handleUpdateNode,
  isUpdateAvailable,
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
            loading={loading}
            isClickable={true}
            selected={selected}
            title={'Node Information'}
            handleAction={handleOnSelected}
          >
            <NodeStatItem
              value={`${capitalize(selectedNode?.type.toLowerCase() ?? 'HOME')} Node`}
              name={'Model type'}
            />

            <NodeStatItem
              value={selectedNode?.id.toLowerCase() ?? '-'}
              name={'Serial #'}
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
            loading={loading}
            isClickable={true}
            selected={selected}
            title={'Node Health'}
            handleAction={handleOnSelected}
          >
            {healthConfig.map((config) => (
              <NodeStatItem
                key={config.id}
                name={config.name}
                value={KPI_PLACEHOLDER_VALUE}
                nameInfo={config.description}
              />
            ))}
          </NodeStatsContainer>
          {selectedNode?.type !== NodeTypeEnum.Anode && (
            <NodeStatsContainer
              index={2}
              loading={loading}
              isClickable={true}
              selected={selected}
              title={'Subscribers'}
              handleAction={handleOnSelected}
            >
              {subscriberConfig.map((config) => (
                <NodeStatItem
                  key={config.id}
                  name={config.name}
                  value={KPI_PLACEHOLDER_VALUE}
                  nameInfo={config.description}
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
            loading={loading}
            nodeTitle={selectedNode?.name ?? 'HOME'}
            handleUpdateNode={handleUpdateNode}
            isUpdateAvailable={isUpdateAvailable}
          />
        )}
        <Paper
          sx={{
            p: 3,
            overflow: 'auto',
            height: { xs: 'calc(100vh - 480px)', md: 'calc(100vh - 328px)' },
          }}
        >
          {selected === 1 && (
            <Stack spacing={4}>
              {healthConfig.map((config) => (
                <LineChart
                  key={config.id}
                  from={metricFrom}
                  loading={metricsLoading}
                  topic={config.id}
                  title={config.description}
                  hasData={isMetricValue(config.id, metrics)}
                  initData={getMetricValue(config.id, metrics)}
                />
              ))}
            </Stack>
          )}
          {selected === 2 && nodeType === NodeTypeEnum.Tnode && (
            <Stack spacing={4}>
              {subscriberConfig.map((config) => (
                <LineChart
                  key={config.id}
                  from={metricFrom}
                  loading={metricsLoading}
                  topic={config.id}
                  title={config.description}
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
