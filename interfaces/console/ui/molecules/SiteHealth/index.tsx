import React from 'react';
import { Box, Stack, Typography, Paper, Grid } from '@mui/material';
import { SiteHealth } from '@/public/svg';

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
    const renderBatteryInfo = React.useCallback((batteryInfo: BatteryInfo) => {
      const { label, value } = batteryInfo;
      return (
        <Box
          display="flex"
          component={Grid}
          flexDirection="row"
          alignItems="center"
          spacing={5}
        >
          <Typography variant="body1" color="initial">
            {label}:
          </Typography>
          <Typography variant="body1" color="initial">
            {value}
          </Typography>
        </Box>
      );
    }, []);

    const batteryCharge = React.useMemo(
      () => `${batteryInfo[2].value} %`,
      [batteryInfo],
    );
    const batteryVoltage = React.useMemo(
      () => `${batteryInfo[0].value} V`,
      [batteryInfo],
    );
    const batteryCurrent = React.useMemo(
      () => `${batteryInfo[1].value} A`,
      [batteryInfo],
    );
    const batteryPower = React.useMemo(
      () => `${batteryInfo[3].value} W`,
      [batteryInfo],
    );

    return (
      <Box>
        <Grid container spacing={2}>
          <Grid item xs={8}>
            <SiteHealth
              solarHealth={solarHealth}
              nodeHealth={nodeHealth}
              switchHealth={switchHealth}
              controllerHealth={controllerHealth}
              batteryHealth={batteryHealth}
              backhaulHealth={backhaulHealth}
            />
          </Grid>
          <Grid item xs={4}>
            <Stack direction={'column'} spacing={2}>
              <Paper variant="outlined" sx={{ p: 2 }}>
                <Stack direction="column" spacing={1} sx={{ p: 1 }}>
                  <Typography variant="h6" color="initial">
                    Battery information
                  </Typography>
                  {batteryInfo.map(renderBatteryInfo)}
                </Stack>
              </Paper>
              <Paper variant="outlined" sx={{ p: 2 }}>
                <Stack direction="column" spacing={1} sx={{ p: 1 }}>
                  <Typography variant="h6" color="initial">
                    Battery information
                  </Typography>
                  {renderBatteryInfo({
                    label: 'Charge',
                    value: batteryCharge,
                  })}
                  {renderBatteryInfo({
                    label: 'Voltage',
                    value: batteryVoltage,
                  })}
                  {renderBatteryInfo({
                    label: 'Current',
                    value: batteryCurrent,
                  })}
                  {renderBatteryInfo({
                    label: 'Power',
                    value: batteryPower,
                  })}
                </Stack>
              </Paper>
            </Stack>
          </Grid>
        </Grid>
      </Box>
    );
  },
);

export default SiteOverallHealth;
SiteOverallHealth.displayName = 'SiteOverallHealth';
