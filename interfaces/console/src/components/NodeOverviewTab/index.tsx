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
import { HealtChartsConfigure, TooltipsText } from '@/constants';
import { getMetricValue, isMetricValue } from '@/utils';
import { Paper, Stack, Typography, capitalize } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { useEffect, useState } from 'react';
import LineChart from '../LineChart';
import NodeDetailsCard from '../NodeDetailsCard';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

const PLACEHOLDER_VALUE = 'NA';

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
  const nodeType = selectedNode?.type ?? NodeTypeEnum.Hnode;
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
      <Grid size={{ xs: 12, md: 4 }}>
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
            {HealtChartsConfigure[nodeType][0].show && (
              <NodeStatItem
                value={PLACEHOLDER_VALUE}
                name={HealtChartsConfigure[nodeType][0].name}
                showAlertInfo={false}
                nameInfo={TooltipsText.TRX}
              />
            )}
            {HealtChartsConfigure[nodeType][1].show && (
              <NodeStatItem
                value={PLACEHOLDER_VALUE}
                name={HealtChartsConfigure[nodeType][1].name}
                nameInfo={TooltipsText.COM}
              />
            )}
            {HealtChartsConfigure[nodeType][2].show && (
              <NodeStatItem
                name={HealtChartsConfigure[nodeType][2].name}
                nameInfo={TooltipsText.COM}
                value={PLACEHOLDER_VALUE}
              />
            )}
          </NodeStatsContainer>
          {/* {selectedNode?.type !== NodeTypeEnum.Anode && (
            <NodeStatsContainer
              index={2}
              loading={loading}
              isClickable={true}
              selected={selected}
              title={'Subscribers'}
              handleAction={handleOnSelected}
            >
              <NodeStatItem
                name={'Attached'}
                value={connectedUsers}
                nameInfo={TooltipsText.ATTACHED}
              />
              <NodeStatItem
                name={'Active'}
                value={`${
                  connectedUsers === '0'
                    ? parseInt(connectedUsers)
                    : parseInt(connectedUsers) - 1
                }`}
                nameInfo={TooltipsText.ACTIVE}
              />
            </NodeStatsContainer>
          )} */}
        </Stack>
      </Grid>
      <Grid size={{ xs: 12, md: 8 }}>
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
        {selected === 1 && (
          <Paper
            sx={{
              p: 3,
              overflow: 'auto',
              height: { xs: 'calc(100vh - 480px)', md: 'calc(100vh - 328px)' },
            }}
          >
            <Stack spacing={4}>
              <Typography variant="h6">Node Health</Typography>
              {HealtChartsConfigure[nodeType][0].show && (
                <LineChart
                  nodeId={nodeId}
                  tabSection={Graphs_Type.NodeHealth}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  topic={HealtChartsConfigure[nodeType][0].id}
                  title={HealtChartsConfigure[nodeType][0].name}
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][0].id,
                    metrics,
                  )}
                  hasData={isMetricValue(
                    HealtChartsConfigure[nodeType][0].id,
                    metrics,
                  )}
                />
              )}
              {HealtChartsConfigure[nodeType][1].show && (
                <LineChart
                  nodeId={nodeId}
                  tabSection={Graphs_Type.NodeHealth}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  topic={HealtChartsConfigure[nodeType][1].id}
                  title={HealtChartsConfigure[nodeType][1].name}
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][1].id,
                    metrics,
                  )}
                  hasData={isMetricValue(
                    HealtChartsConfigure[nodeType][1].id,
                    metrics,
                  )}
                />
              )}
              {HealtChartsConfigure[nodeType][2].show && (
                <LineChart
                  nodeId={nodeId}
                  tabSection={Graphs_Type.NodeHealth}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  topic={HealtChartsConfigure[nodeType][2].id}
                  title={HealtChartsConfigure[nodeType][2].name}
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][2].id,
                    metrics,
                  )}
                  hasData={isMetricValue(
                    HealtChartsConfigure[nodeType][2].id,
                    metrics,
                  )}
                />
              )}
            </Stack>
          </Paper>
        )}
        {/* {selected === 2 && nodeType !== NodeTypeEnum.Anode && (
          <Paper sx={{ p: 3 }}>
            <Stack spacing={4}>
              <Typography variant="h6">Subscribers</Typography>
              {HealtChartsConfigure[
                (selectedNode?.type as string) ?? 'hnode'
              ][4].show && (
                <LineChart
                  nodeId={nodeId}
                  tabSection={Graphs_Type.Subscribers}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  topic={HealtChartsConfigure[nodeType][3].id}
                  title={HealtChartsConfigure[nodeType][3].name}
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][3].id,
                    metrics,
                  )}
                  hasData={isMetricValue(
                    HealtChartsConfigure[nodeType][3].id,
                    metrics,
                  )}
                />
              )}
              {HealtChartsConfigure[
                (selectedNode?.type as string) ?? 'hnode'
              ][4].show && (
                <LineChart
                  nodeId={nodeId}
                  tabSection={Graphs_Type.Subscribers}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  topic={HealtChartsConfigure[nodeType][4].id}
                  title={HealtChartsConfigure[nodeType][4].name}
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][4].id,
                    metrics,
                  )}
                  hasData={isMetricValue(
                    HealtChartsConfigure[nodeType][4].id,
                    metrics,
                  )}
                />
              )}
            </Stack>
          </Paper>
        )} */}
      </Grid>
    </Grid>
  );
};

export default NodeOverviewTab;
