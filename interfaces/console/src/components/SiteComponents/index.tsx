/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Card, Grid, Stack, Paper } from '@mui/material';
import LineChart from '../LineChart';
import { MetricsRes } from '@/client/graphql/generated/subscriptions';
import { getMetricValue, isMetricValue } from '@/utils';
import LoadingWrapper from '@/components/LoadingWrapper';
import SiteFlowDiagram from '../../../public/svg/sitecomps';
import NodeStatusDisplay from '@/components/NodeStatusDisplay';

export interface SiteKpiConfig {
  id: string;
  name: string;
  unit: string;
  description: string;
  tickInterval?: number;
  tickPositions?: number[];
  threshold?: {
    min: number;
    normal: number;
    max: number;
  } | null;
  show?: boolean;
}

export interface SectionData {
  [key: string]: SiteKpiConfig[];
}

export const KPI_TO_SECTION_MAP: Record<string, string> = {
  solar: 'SOLAR',
  battery: 'BATTERY',
  controller: 'CONTROLLER',
  backhaul: 'MAIN_BACKHAUL',
  switch: 'SWITCH',
  node: 'NODE',
};

interface SiteComponentsProps {
  siteId: string;
  metrics: MetricsRes;
  sections: SectionData;
  nodeIds: string[];
  activeKPI: string;
  activeSection: string;
  metricFrom: number;
  metricsLoading: boolean;
  onNodeClick: (kpiType: string) => void;
}

const SiteComponents: React.FC<SiteComponentsProps> = ({
  metrics,
  sections,
  nodeIds,
  activeKPI,
  activeSection,
  metricFrom,
  metricsLoading,
  onNodeClick,
}) => {
  return (
    <Box>
      <Card
        sx={{
          p: 3,
          borderRadius: 2,
          boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.05)',
          width: '100%',
        }}
      >
        <Grid container spacing={3}>
          <Grid item xs={12} md={3}>
            <Stack spacing={2}>
              <SiteFlowDiagram onNodeClick={onNodeClick} />
            </Stack>
          </Grid>

          <Grid item xs={12} md={9}>
            <LoadingWrapper
              radius="small"
              width="100%"
              isLoading={metricsLoading && activeKPI !== 'node'}
            >
              {activeKPI === 'node' ? (
                <NodeStatusDisplay nodeIds={nodeIds} />
              ) : (
                <Paper
                  sx={{
                    p: 3,
                    overflow: 'auto',
                    height: {
                      xs: 'calc(100vh - 480px)',
                      md: 'calc(100vh - 328px)',
                    },
                  }}
                >
                  <Stack spacing={4}>
                    {sections[activeSection]?.map((config, i) => (
                      <LineChart
                        from={metricFrom}
                        topic={config.id}
                        title={config.name}
                        yunit={config.unit}
                        loading={metricsLoading}
                        key={`${config.id}-${i}`}
                        tickInterval={config.tickInterval}
                        tickPositions={config.tickPositions}
                        hasData={isMetricValue(config.id, metrics)}
                        initData={getMetricValue(config.id, metrics)}
                      />
                    ))}
                  </Stack>
                </Paper>
              )}
            </LoadingWrapper>
          </Grid>
        </Grid>
      </Card>
    </Box>
  );
};

export default SiteComponents;
