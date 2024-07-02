/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useState } from 'react';
import { Grid, Box, Typography } from '@mui/material';
import { SiteHealth } from '@/../public/svg';
import colors from '@/theme/colors';

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
                {renderKpiInfo()}{' '}
              </Typography>
            </Box>
          </Grid>
        </Grid>
      </>
    );
  },
);

SiteOverallHealth.displayName = 'SiteOverallHealth';

export default SiteOverallHealth;
