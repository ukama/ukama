/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  Grid,
  Stack,
  Paper,
  Typography,
  Switch,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import LineChart from '../LineChart';
import { MetricsRes } from '@/client/graphql/generated/subscriptions';
import { getMetricValue, getPortInfo, isMetricValue } from '@/utils';
import SiteFlowDiagram from '../../../public/svg/sitecomps';
import NodeStatusDisplay from '@/components/NodeStatusDisplay';
import { SectionData, SITE_KPI_TYPES } from '@/constants';

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
  nodeIds?: string[];
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
  nodeIds,
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
    if (!siteId || !nodeIds || nodeIds.length === 0) return;

    const tokens = nodeIds.map((nodeId) => {
      const topic = `stat-${SITE_KPI_TYPES.NODE_UPTIME}-${nodeId}`;
      return PubSub.subscribe(topic, (_, uptimeValue) => {
        setNodeUptimes((prev) => ({
          ...prev,
          [nodeId]: Math.floor(uptimeValue[1]),
        }));
      });
    });

    return () => {
      [...tokens].forEach((token) => PubSub.unsubscribe(token));
    };
  }, [siteId, nodeIds]);

  useEffect(() => {
    if (hasMetricsData && activeSection === 'SWITCH') {
      const portGroups = getPortMetrics();

      const newSwitchStatus: Record<string, boolean> = {};

      portGroups.forEach((portGroup) => {
        const statusMetric = portGroup.metrics.find((m: any) =>
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
  }, [metrics, activeSection, hasMetricsData]);

  const togglePortExpand = (portId: string) => {
    setExpandedPorts((prev) => ({
      ...prev,
      [portId]: !prev[portId],
    }));
  };

  const getPortMetrics = () => {
    const switchMetrics = sections[activeSection] || [];

    const portGroups: Record<string, any[]> = {};

    switchMetrics.forEach((metric) => {
      let portType = '';
      if (metric.id.startsWith('solar_')) portType = 'solar';
      else if (metric.id.startsWith('backhaul_')) portType = 'backhaul';
      else if (metric.id.startsWith('node_')) portType = 'node';
      else return;

      if (!portGroups[portType]) {
        portGroups[portType] = [];
      }

      portGroups[portType].push(metric);
    });

    return Object.entries(portGroups)
      .map(([portType, metrics]) => {
        const port = getPortInfo[portType] || { number: 0, desc: '' };
        return {
          id: portType,
          portNumber: port.number,
          description: port.desc,
          metrics,
        };
      })
      .sort((a, b) => a.portNumber - b.portNumber);
  };

  const renderPortItem = (portGroup: any) => {
    const isExpanded = expandedPorts[portGroup.id] || false;
    const portTitle = portGroup.description
      ? `Port ${portGroup.portNumber} (${portGroup.description})`
      : `Port ${portGroup.portNumber}`;

    const statusMetric = portGroup.metrics.find((m: any) =>
      m.id.includes('switch_port_status'),
    );

    const isOn = localSwitchStatus[portGroup.id] ?? false;
    const isDisabled = disabledSwitches[portGroup.id] || false;

    const handleToggle = () => {
      if (onSwitchChange) {
        setLocalSwitchStatus((prev) => ({
          ...prev,
          [portGroup.id]: !isOn,
        }));

        setDisabledSwitches((prev) => ({
          ...prev,
          [portGroup.id]: true,
        }));

        onSwitchChange(portGroup.portNumber, isOn);

        setTimeout(() => {
          setDisabledSwitches((prev) => ({
            ...prev,
            [portGroup.id]: false,
          }));
        }, 5000);
      }
    };

    return (
      <Box
        key={portGroup.id}
        sx={{ borderBottom: '1px solid rgba(0, 0, 0, 0.12)', py: 2 }}
      >
        <Accordion
          expanded={isExpanded}
          onChange={() => togglePortExpand(portGroup.id)}
          sx={{ boxShadow: 'none' }}
        >
          <AccordionSummary
            expandIcon={<ExpandMoreIcon />}
            sx={{ display: 'flex', alignItems: 'center', p: 0 }}
          >
            <Typography
              variant="subtitle1"
              fontWeight="medium"
              sx={{ flexGrow: 1 }}
            >
              {portTitle}
            </Typography>
          </AccordionSummary>
          <AccordionDetails sx={{ mt: 2, ml: 2, p: 0 }}>
            {statusMetric && (
              <Box
                sx={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  mb: 3,
                  pb: 2,
                  borderBottom: '1px solid rgba(0, 0, 0, 0.06)',
                }}
              >
                <Typography variant="body1">{statusMetric.name}</Typography>
                <Box display="flex" alignItems="center">
                  <Typography variant="body2" sx={{ mr: 1 }}>
                    {isOn ? 'On' : 'Off'}
                  </Typography>
                  <Switch
                    checked={isOn}
                    onChange={handleToggle}
                    disabled={isDisabled}
                    color="primary"
                  />
                </Box>
              </Box>
            )}

            <Stack spacing={3}>
              {portGroup.metrics
                .filter((m: any) => !m.id.includes('switch_port_status'))
                .map((metric: any) => (
                  <Box key={metric.id}>
                    <Typography variant="body1" sx={{ mb: 1 }}>
                      {metric.name}
                    </Typography>
                    <LineChart
                      from={metricFrom}
                      topic={metric.id}
                      yunit={metric.unit}
                      loading={metricsLoading}
                      tickInterval={metric.tickInterval}
                      tickPositions={metric.tickPositions}
                      hasData={isMetricValue(metric.id, metrics)}
                      initData={getMetricValue(metric.id, metrics)}
                      format={metric.format}
                    />
                  </Box>
                ))}
            </Stack>
          </AccordionDetails>
        </Accordion>
      </Box>
    );
  };

  return (
    <Box>
      <Card
        sx={{
          p: 3,
          borderRadius: 2,
          width: '100%',
        }}
      >
        <Grid container spacing={3}>
          <Grid item xs={12} md={3}>
            <Stack spacing={2}>
              <SiteFlowDiagram
                defaultOpacity={0.1}
                onNodeClick={onComponentClick}
                activeKPI={activeKPI}
              />
            </Stack>
          </Grid>

          <Grid item xs={12} md={9}>
            {activeKPI === 'node' ? (
              <NodeStatusDisplay nodeUptimes={nodeUptimes} />
            ) : (
              <Paper
                sx={{
                  p: 3,
                  pr: 5,
                  overflow: 'auto',
                  height: {
                    xs: 'calc(100vh - 480px)',
                    md: 'calc(100vh - 328px)',
                  },
                }}
              >
                {activeSection === 'SWITCH' && (
                  <Box sx={{ mb: 4 }}>
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

                      {getPortMetrics().map((portGroup) =>
                        renderPortItem(portGroup),
                      )}
                    </Box>
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
      </Card>
    </Box>
  );
};

export default SiteComponents;
