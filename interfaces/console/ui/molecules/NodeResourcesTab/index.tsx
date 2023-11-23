/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NodeResourcesTabConfigure, TooltipsText } from '@/constants';
import { Node, NodeTypeEnum } from '@/generated';
import { Graphs_Type } from '@/generated/metrics';
import { getMetricValue, isMetricValue } from '@/utils';
import { Grid, Paper, Stack } from '@mui/material';
import { useState } from 'react';
import LineChart from '../LineChart';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

const PLACEHOLDER_VALUE = 'NA';
interface INodeResourcesTab {
  metrics: any;
  loading: boolean;
  metricFrom: number;
  selectedNode: Node | undefined;
}
const NodeResourcesTab = ({
  metrics,
  loading,
  metricFrom,
  selectedNode,
}: INodeResourcesTab) => {
  const nodeType = selectedNode?.type || NodeTypeEnum.Hnode;
  const [isCollapse, setIsCollapse] = useState<boolean>(false);
  const handleCollapse = () => setIsCollapse((prev) => !prev);
  return (
    <Grid container spacing={3}>
      <Grid item lg={!isCollapse ? 4 : 1} md xs>
        <NodeStatsContainer
          index={0}
          selected={0}
          loading={loading}
          title={'Resources'}
          isCollapsable={true}
          isCollapse={isCollapse}
          onCollapse={handleCollapse}
        >
          {NodeResourcesTabConfigure[nodeType][0].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              variant={'large'}
              name={NodeResourcesTabConfigure[nodeType][0].name}
              nameInfo={TooltipsText.MTRX}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][1].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              variant={'large'}
              name={NodeResourcesTabConfigure[nodeType][1].name}
              nameInfo={TooltipsText.MCOM}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][2].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              name={NodeResourcesTabConfigure[nodeType][2].name}
              variant={'large'}
              nameInfo={TooltipsText.CPUTRX}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][3].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              name={NodeResourcesTabConfigure[nodeType][3].name}
              variant={'large'}
              nameInfo={TooltipsText.CPUCOM}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][4].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              variant={'large'}
              name={NodeResourcesTabConfigure[nodeType][4].name}
              nameInfo={TooltipsText.DISKTRX}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][5].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              variant={'large'}
              name={NodeResourcesTabConfigure[nodeType][5].name}
              nameInfo={TooltipsText.DISKCOM}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][6].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              name={NodeResourcesTabConfigure[nodeType][6].name}
              variant={'large'}
              nameInfo={TooltipsText.POWER}
            />
          )}
        </NodeStatsContainer>
      </Grid>
      <Grid item lg={isCollapse ? 11 : 8} md xs>
        <Paper sx={{ p: 3, width: '100%' }}>
          <Stack spacing={4}>
            {NodeResourcesTabConfigure[nodeType][0].show && (
              <LineChart
                loading={loading}
                initData={getMetricValue(
                  NodeResourcesTabConfigure[nodeType][0].id,
                  metrics,
                )}
                metricFrom={metricFrom}
                tabSection={Graphs_Type.Resources}
                topic={NodeResourcesTabConfigure[nodeType][0].id}
                title={NodeResourcesTabConfigure[nodeType][0].name}
                hasData={isMetricValue(
                  NodeResourcesTabConfigure[nodeType][0].id,
                  metrics,
                )}
              />
            )}
            {NodeResourcesTabConfigure[nodeType][1].show && (
              <LineChart
                loading={loading}
                initData={getMetricValue(
                  NodeResourcesTabConfigure[nodeType][1].id,
                  metrics,
                )}
                metricFrom={metricFrom}
                tabSection={Graphs_Type.Resources}
                topic={NodeResourcesTabConfigure[nodeType][1].id}
                title={NodeResourcesTabConfigure[nodeType][1].name}
                hasData={isMetricValue(
                  NodeResourcesTabConfigure[nodeType][1].id,
                  metrics,
                )}
              />
            )}
            {NodeResourcesTabConfigure[nodeType][2].show && (
              <LineChart
                loading={loading}
                initData={getMetricValue(
                  NodeResourcesTabConfigure[nodeType][2].id,
                  metrics,
                )}
                metricFrom={metricFrom}
                tabSection={Graphs_Type.Resources}
                topic={NodeResourcesTabConfigure[nodeType][2].id}
                title={NodeResourcesTabConfigure[nodeType][2].name}
                hasData={isMetricValue(
                  NodeResourcesTabConfigure[nodeType][2].id,
                  metrics,
                )}
              />
            )}
            {NodeResourcesTabConfigure[nodeType][3].show && (
              <LineChart
                loading={loading}
                initData={getMetricValue(
                  NodeResourcesTabConfigure[nodeType][3].id,
                  metrics,
                )}
                metricFrom={metricFrom}
                tabSection={Graphs_Type.Resources}
                topic={NodeResourcesTabConfigure[nodeType][3].id}
                title={NodeResourcesTabConfigure[nodeType][3].name}
                hasData={isMetricValue(
                  NodeResourcesTabConfigure[nodeType][3].id,
                  metrics,
                )}
              />
            )}
            {NodeResourcesTabConfigure[nodeType][4].show && (
              <LineChart
                loading={loading}
                initData={getMetricValue(
                  NodeResourcesTabConfigure[nodeType][4].id,
                  metrics,
                )}
                metricFrom={metricFrom}
                tabSection={Graphs_Type.Resources}
                topic={NodeResourcesTabConfigure[nodeType][4].id}
                title={NodeResourcesTabConfigure[nodeType][4].name}
                hasData={isMetricValue(
                  NodeResourcesTabConfigure[nodeType][4].id,
                  metrics,
                )}
              />
            )}
            {NodeResourcesTabConfigure[nodeType][5].show && (
              <LineChart
                loading={loading}
                initData={getMetricValue(
                  NodeResourcesTabConfigure[nodeType][5].id,
                  metrics,
                )}
                metricFrom={metricFrom}
                tabSection={Graphs_Type.Resources}
                topic={NodeResourcesTabConfigure[nodeType][5].id}
                title={NodeResourcesTabConfigure[nodeType][5].name}
                hasData={isMetricValue(
                  NodeResourcesTabConfigure[nodeType][5].id,
                  metrics,
                )}
              />
            )}
            {NodeResourcesTabConfigure[nodeType][6].show && (
              <LineChart
                loading={loading}
                initData={getMetricValue(
                  NodeResourcesTabConfigure[nodeType][6].id,
                  metrics,
                )}
                metricFrom={metricFrom}
                tabSection={Graphs_Type.Resources}
                topic={NodeResourcesTabConfigure[nodeType][6].id}
                title={NodeResourcesTabConfigure[nodeType][6].name}
                hasData={isMetricValue(
                  NodeResourcesTabConfigure[nodeType][6].id,
                  metrics,
                )}
              />
            )}
          </Stack>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeResourcesTab;
