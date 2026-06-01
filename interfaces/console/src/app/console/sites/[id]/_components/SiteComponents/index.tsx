/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import NodeStatusDisplay from '@/app/console/nodes/[id]/_components/NodeStatusDisplay';
import { Node } from '@/client/graphql/generated';
import { MetricsRes } from '@/client/graphql/generated/subscriptions';
import { SectionData, SITE_KPI_TYPES, SiteKpiConfig } from '@/constants';
import { getMetricValue, isMetricValue } from '@/utils';
// getMetricValue / isMetricValue are still used for non-SWITCH metric charts below
import LineChart from '@/components/ui/LineChart';
import { Box, Paper, Stack, Typography } from '@mui/material';
import Grid from '@mui/material/Grid2';
import React, { useEffect, useState } from 'react';
import SiteFlowDiagram from '../../../../../../../public/svg/sitecomps';
import SwitchPortItem, { PortGroup } from './SwitchPortItem';
interface SiteComponentsProps {
  siteId: string;
  metrics: MetricsRes;
  sections: SectionData;
  activeKPI: string;
  activeSection: string;
  metricFrom: number;
  metricsLoading: boolean;
  onComponentClick: (kpiType: string) => void;
  onSwitchChange?: (portNumber: number, currentStatus: boolean) => void;
  nodes?: Node[];
  initialNodeUptimes?: Record<string, number>;
}

const SiteComponents: React.FC<SiteComponentsProps> = ({
  metrics,
  sections,
  activeKPI,
  activeSection,
  siteId,
  metricFrom,
  metricsLoading,
  onComponentClick,
  onSwitchChange,
  nodes,
  initialNodeUptimes,
}) => {
  const hasMetricsData =
    metrics && metrics.metrics && metrics.metrics.length > 0;

  const [expandedPorts, setExpandedPorts] = useState<Record<string, boolean>>(
    {},
  );

  const [disabledSwitches, setDisabledSwitches] = useState<
    Record<string, boolean>
  >({});

  const [localSwitchStatus, setLocalSwitchStatus] = useState<
    Record<string, boolean>
  >({});

  const [nodeUptimes, setNodeUptimes] = useState<Record<string, number>>({});
  useEffect(() => {
    if (initialNodeUptimes && Object.keys(initialNodeUptimes).length > 0) {
      setNodeUptimes(initialNodeUptimes);
    }
  }, [initialNodeUptimes]);
  useEffect(() => {
    if (!siteId || !nodes || nodes.length === 0) return;

    const tokens = nodes.map((node) => {
      const topic = `stat-${SITE_KPI_TYPES.NODE_UPTIME}-${node.id}`;
      return PubSub.subscribe(topic, (_, uptimeValue) => {
        setNodeUptimes((prev) => ({
          ...prev,
          [node.id]: Math.floor(uptimeValue[1]),
        }));
      });
    });

    return () => {
      [...tokens].forEach((token) => PubSub.unsubscribe(token));
    };
  }, [siteId, nodes]);

  useEffect(() => {
    if (hasMetricsData && activeSection === 'SWITCH') {
      const portGroups = getPortMetrics();

      const newSwitchStatus: Record<string, boolean> = {};

      portGroups.forEach((portGroup) => {
        const statusMetric = portGroup.metrics.find((m: SiteKpiConfig) =>
          m.id.includes('switch_port_status'),
        );

        if (statusMetric) {
          const metricValues = getMetricValue(statusMetric.id, metrics);
          if (metricValues && metricValues.length > 0) {
            const latestValue = metricValues[metricValues.length - 1];
            let isOn = false;

            if (Array.isArray(latestValue)) {
              isOn = latestValue[1] === 1;
            } else {
              isOn = latestValue === 1;
            }

            newSwitchStatus[portGroup.id] = isOn;
          }
        }
      });

      if (Object.keys(newSwitchStatus).length > 0) {
        setLocalSwitchStatus((prev) => ({
          ...prev,
          ...newSwitchStatus,
        }));
      }
    }
  }, [metrics, activeSection, hasMetricsData, sections]);

  const togglePortExpand = (portId: string) => {
    setExpandedPorts((prev) => ({
      ...prev,
      [portId]: !prev[portId],
    }));
  };

  const SWITCH_PORT_DESCRIPTIONS: Record<number, string> = {
    1: 'Tower node',
    2: 'Amplifier node',
    3: 'Controller node',
    9: 'Backhaul node',
  };

  const resolveSwitchPortNumber = (metric: SiteKpiConfig): number | null => {
    if (typeof metric.port === 'number' && !Number.isNaN(metric.port)) {
      return metric.port;
    }
    const m = /^switch_port_(\d+)_/.exec(metric.id);
    if (m) return Number.parseInt(m[1], 10);
    return null;
  };

  const getPortMetrics = () => {
    const switchMetrics = sections[activeSection] || [];

    const byPort: Record<number, SiteKpiConfig[]> = {};

    switchMetrics.forEach((metric) => {
      const portNum = resolveSwitchPortNumber(metric);
      if (portNum == null) return;
      if (!byPort[portNum]) byPort[portNum] = [];
      byPort[portNum].push(metric);
    });

    return Object.entries(byPort)
      .map(([portStr, portMetrics]) => {
        const portNumber = Number.parseInt(portStr, 10);
        const description = SWITCH_PORT_DESCRIPTIONS[portNumber] ?? '';
        return {
          id: `port-${portNumber}`,
          portNumber,
          description,
          metrics: portMetrics,
        };
      })
      .sort((a, b) => a.portNumber - b.portNumber);
  };

  const handlePortToggleSwitch = (portNumber: number, currentIsOn: boolean) => {
    if (!onSwitchChange) return;
    const portId = `port-${portNumber}`;
    setLocalSwitchStatus((prev) => ({ ...prev, [portId]: !currentIsOn }));
    setDisabledSwitches((prev) => ({ ...prev, [portId]: true }));
    onSwitchChange(portNumber, currentIsOn);
    setTimeout(() => {
      setDisabledSwitches((prev) => ({ ...prev, [portId]: false }));
    }, 5000);
  };

  return (
    <Paper
      sx={{
        p: 2,
        borderRadius: 2,
        minHeight: 'fit-content',
        height: 'fit-content',
      }}
    >
      <Grid container spacing={3}>
        <Grid size={{ xs: 12, md: 3 }} alignSelf="center">
          <SiteFlowDiagram
            defaultOpacity={0.1}
            onNodeClick={onComponentClick}
            activeKPI={activeKPI}
          />
        </Grid>

        <Grid size={{ xs: 12, md: 9 }}>
          {activeKPI === 'node' ? (
            <NodeStatusDisplay nodes={nodes ?? []} nodeUptimes={nodeUptimes} />
          ) : (
            <Paper
              elevation={0}
              sx={{
                p: 3,
                pr: 5,
                boxShadow: 'none',
                overflow: 'auto',
                height: {
                  xs: 'calc(100vh - 480px)',
                  md: 'calc(100vh - 328px)',
                },
              }}
            >
              {activeSection === 'SWITCH' && (
                <Box
                  sx={{
                    p: 2,
                    borderRadius: 1,
                  }}
                >
                  <Typography
                    variant="h6"
                    fontWeight="medium"
                    sx={{
                      mb: 2,
                      pb: 1,
                      borderBottom: '1px solid rgba(0, 0, 0, 0.12)',
                    }}
                  >
                    {(() => {
                      const portMetrics = getPortMetrics();
                      const totalPorts = portMetrics.length;
                      const activePorts = Object.values(
                        localSwitchStatus,
                      ).filter((status) => status === true).length;
                      return `Switch ports (${activePorts} active / ${totalPorts} total)`;
                    })()}
                  </Typography>

                  {getPortMetrics().map((portGroup: PortGroup) => (
                    <SwitchPortItem
                      key={portGroup.id}
                      portGroup={portGroup}
                      metrics={metrics}
                      metricFrom={metricFrom}
                      metricsLoading={metricsLoading}
                      isExpanded={expandedPorts[portGroup.id] ?? false}
                      isOn={localSwitchStatus[portGroup.id] ?? false}
                      isDisabled={disabledSwitches[portGroup.id] ?? false}
                      onToggleExpand={togglePortExpand}
                      onToggleSwitch={handlePortToggleSwitch}
                    />
                  ))}
                </Box>
              )}

              {activeSection !== 'SWITCH' && (
                <Stack spacing={4}>
                  {sections[activeSection]?.map((config) => (
                    <LineChart
                      from={metricFrom}
                      topic={config.id}
                      title={config.name}
                      yunit={config.unit}
                      loading={metricsLoading}
                      key={config.id}
                      tickInterval={config.tickInterval}
                      tickPositions={config.tickPositions}
                      hasData={isMetricValue(config.id, metrics)}
                      initData={getMetricValue(config.id, metrics)}
                      format={config.format}
                    />
                  ))}
                </Stack>
              )}
            </Paper>
          )}
        </Grid>
      </Grid>
    </Paper>
  );
};

export default SiteComponents;
