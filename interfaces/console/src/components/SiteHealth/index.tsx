import React, { useState, useEffect } from 'react';
import { Grid, Box, Typography, Stack } from '@mui/material';
import LineChart from '../LineChart';
import { getMetricValue, isMetricValue } from '@/utils';
import {
  Graphs_Type,
  MetricsRes,
} from '@/client/graphql/generated/subscriptions';

interface SiteOverallHealthProps {
  nodeId: string;
  metricFrom: number;
  metrics: MetricsRes;
  loading?: boolean;
  tabSection?: Graphs_Type;
  onSiteKpiChange?: (kpi: Graphs_Type) => void;
}

const POWER_METRICS = [
  {
    id: 'power_input',
    name: 'Power Input',
    show: true,
  },
  {
    id: 'solar_performance',
    name: 'Solar Performance',
    show: true,
  },
  {
    id: 'battery_status',
    name: 'Battery Status',
    show: true,
  },
];

const SiteOverallHealth: React.FC<SiteOverallHealthProps> = React.memo(
  ({
    nodeId,
    metricFrom,
    metrics,
    loading = false,
    tabSection = Graphs_Type.Power,
    onSiteKpiChange,
  }) => {
    // Add debug logs
    useEffect(() => {
      console.log('SiteHealth received props:', {
        nodeId,
        metricFrom,
        metrics,
        loading,
        tabSection,
      });
    }, [nodeId, metricFrom, metrics, loading, tabSection]);

    const [selectedKpi, setSelectedKpi] = useState<Graphs_Type>(tabSection);

    useEffect(() => {
      // Initialize with Power metrics
      onSiteKpiChange?.(Graphs_Type.Power);
      setSelectedKpi(Graphs_Type.Power);
    }, [nodeId, onSiteKpiChange]);

    return (
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Typography variant="h6" sx={{ mb: 3 }}>
            Site Power Metrics{' '}
            {nodeId ? `(Node: ${nodeId})` : '(No node selected)'}
          </Typography>
          <Stack spacing={4}>
            {POWER_METRICS.map(
              (metric) =>
                metric.show && (
                  <Box key={metric.id}>
                    <LineChart
                      nodeId={nodeId}
                      tabSection={selectedKpi}
                      metricFrom={metricFrom}
                      loading={loading}
                      topic={metric.id}
                      title={metric.name}
                      initData={getMetricValue(metric.id, metrics)}
                      hasData={isMetricValue(metric.id, metrics)}
                    />
                  </Box>
                ),
            )}
          </Stack>
        </Grid>
      </Grid>
    );
  },
);

SiteOverallHealth.displayName = 'SiteOverallHealth';

export default SiteOverallHealth;
