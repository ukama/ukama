/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Box, Card, Grid, Stack, Paper, Skeleton } from '@mui/material';
import LineChart from '../LineChart';
import { MetricsRes } from '@/client/graphql/generated/subscriptions';
import { getMetricValue, isMetricValue } from '@/utils';
import LoadingWrapper from '@/components/LoadingWrapper';
import SiteFlowDiagram from '../../../public/svg/sitecomps';
import NodeStatusDisplay from '@/components/NodeStatusDisplay';
import { SectionData } from '@/constants';

interface SiteComponentsProps {
  siteId: string;
  metrics: MetricsRes;
  sections: SectionData;
  nodeIds: string[];
  activeKPI: string;
  nodeUpTime?: number;
  activeSection: string;
  metricFrom: number;
  metricsLoading: boolean;
  onComponentClick: (kpiType: string) => void;
  onSwitchChange?: () => void;
}

const SiteComponents: React.FC<SiteComponentsProps> = ({
  metrics,
  sections,
  nodeIds,
  activeKPI,
  activeSection,
  metricFrom,
  metricsLoading,
  onComponentClick,
  nodeUpTime,
  onSwitchChange,
}) => {
  const hasMetricsData =
    metrics && metrics.metrics && metrics.metrics.length > 0;

  const renderSkeletonLoading = () => {
    return (
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
          {Array.from({ length: sections[activeSection]?.length || 3 }).map(
            (_, index) => (
              <Box key={`skeleton-${index}`}>
                <Skeleton
                  variant="text"
                  width="40%"
                  height={30}
                  sx={{ mb: 1 }}
                />
                <Skeleton variant="rectangular" width="100%" height={200} />
              </Box>
            ),
          )}
        </Stack>
      </Paper>
    );
  };

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
              <SiteFlowDiagram
                defaultOpacity={0.1}
                onNodeClick={onComponentClick}
              />
            </Stack>
          </Grid>

          <Grid item xs={12} md={9}>
            <LoadingWrapper
              radius="small"
              width="100%"
              isLoading={metricsLoading && activeKPI !== 'node'}
            >
              {activeKPI === 'node' ? (
                <NodeStatusDisplay nodeIds={nodeIds} nodeUpTime={nodeUpTime} />
              ) : !hasMetricsData && !metricsLoading ? (
                renderSkeletonLoading()
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
                        switchLabel={
                          config.id === 'switch_port_status'
                            ? 'Backhaul'
                            : undefined
                        }
                        initialSwitchState={
                          config.id === 'switch_port_status' ? true : undefined
                        }
                        onSwitchChange={onSwitchChange}
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
