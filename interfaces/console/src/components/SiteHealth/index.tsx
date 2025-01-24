import React, { useState } from 'react';
import { Grid, Box, Typography } from '@mui/material';
import { SiteHealth } from '@/../public/svg';
import colors from '@/theme/colors';
import LineChart from '../LineChart';
import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import { getMetricValue } from '@/utils';

interface BatteryInfo {
  label: string;
  value: string;
}

interface SiteOverallHealthProps {
  batteryInfo: BatteryInfo[];
  solarHealth: 'good' | 'warning';
  nodeHealth: 'good' | 'warning';
  switchHealth: 'good' | 'warning';
  controllerHealth: 'good' | 'warning';
  batteryHealth: 'good' | 'warning';
  backhaulHealth: 'good' | 'warning';
  nodes: any[];
  nodeId: string;
  metricFrom: number;
  topic: string;
  initData: any;
  title?: string;
  filter?: string;
  hasData?: boolean;
  loading?: boolean;
  tabSection: Graphs_Type;
}

const SiteOverallHealth: React.FC<SiteOverallHealthProps> = React.memo(
  ({
    batteryInfo,
    solarHealth,
    nodeHealth,
    switchHealth,
    controllerHealth,
    batteryHealth,
    backhaulHealth,
    nodes,
    nodeId,
    metricFrom,
    topic,
    initData,
    title = '',
    loading = false,
    filter = 'LIVE',
    tabSection = Graphs_Type.NodeHealth,
  }) => {
    const [selectedKpi, setSelectedKpi] = useState<string | null>(null);

    const handleNodeClick = () => {
      setSelectedKpi('Node');
    };

    const handleSolarClick = () => {
      setSelectedKpi('Solar');
    };

    const handleSwitchClick = () => {
      setSelectedKpi('Switch');
    };

    const handleControllerClick = () => {
      setSelectedKpi('Controller');
    };

    const handleBatteryClick = () => {
      setSelectedKpi('Battery');
    };

    const handleBackhaulClick = () => {
      setSelectedKpi('Backhaul');
    };

    const renderKpiInfo = () => {
      switch (selectedKpi) {
        case 'Node':
          return 'Nodes';
        case 'Solar':
          return 'Solar panels KPIs';
        case 'Switch':
          return 'Switch overview';
        case 'Controller':
          return 'Charge controller overview';
        case 'Battery':
          return 'Batteries KPIs';
        case 'Backhaul':
          return 'Backhaul overview';
        default:
          return 'Please select a KPI to view information';
      }
    };

    const renderNodeInfo = () => {
      return nodes.map((n, index) => (
        <Typography key={index} variant="body1" color="initial">
          Node #{n.id}: {n.status.state} and{' '}
          {n.status.connectivity.toLowerCase()}
        </Typography>
      ));
    };

    return (
      <>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Typography
              variant="body1"
              color="initial"
              sx={{ fontWeight: 'bold' }}
            >
              Site components
            </Typography>
          </Grid>
          <Grid item xs={6}>
            <SiteHealth
              solarHealth={solarHealth}
              nodeHealth={nodeHealth}
              switchHealth={switchHealth}
              controllerHealth={controllerHealth}
              batteryHealth={batteryHealth}
              backhaulHealth={backhaulHealth}
              onNodeClick={handleNodeClick}
              onSolarClick={handleSolarClick}
              onSwitchClick={handleSwitchClick}
              onControllerClick={handleControllerClick}
              onBatteryClick={handleBatteryClick}
              onBackhaulClick={handleBackhaulClick}
            />
          </Grid>
          <Grid item xs={6}>
            <Box sx={{ border: `1px solid ${colors.black40}`, p: 2 }}>
              <Typography
                variant="body1"
                color="initial"
                sx={{ fontWeight: 'bold' }}
              >
                {renderKpiInfo()}
              </Typography>
              {selectedKpi === 'Node' && nodes ? (
                renderNodeInfo()
              ) : (
                <LineChart
                  nodeId={nodeId}
                  tabSection={Graphs_Type.NodeHealth}
                  metricFrom={metricFrom}
                  loading={loading}
                  topic={topic}
                  title={title}
                  initData={initData}
                  hasData={false}
                />
              )}
            </Box>
          </Grid>
        </Grid>
      </>
    );
  },
);

SiteOverallHealth.displayName = 'SiteOverallHealth';

export default SiteOverallHealth;
