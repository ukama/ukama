/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, NodeTypeEnum } from '@/client/graphql/generated';
import { KPI_PLACEHOLDER_VALUE, NODE_KPIS } from '@/constants';
import { getMetricValue, isMetricValue } from '@/utils';
import { Paper, Stack } from '@mui/material';
import Grid from '@mui/material/Grid2';
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
  const resourcesConfig = NODE_KPIS.RESOURCES[NodeTypeEnum.Tnode];

  return (
    <Grid container spacing={3}>
      <Grid size={{ xs: 12, md: 3 }}>
        <NodeStatsContainer
          index={0}
          selected={0}
          loading={loading}
          title={'Resources'}
          isCollapsable={false}
        >
          {resourcesConfig.map((config) => (
            <NodeStatItem
              key={config.id}
              name={config.name}
              value={KPI_PLACEHOLDER_VALUE}
              nameInfo={config.description}
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
            {resourcesConfig.map((config) => (
              <LineChart
                key={config.id}
                from={metricFrom}
                topic={config.id}
                loading={loading}
                title={config.name}
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
