import React, { PureComponent } from 'react';
import { Stack, Typography, Paper, Grid } from '@mui/material';
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

class SiteOverallHealth extends PureComponent<SiteOverallHealthProps> {
  renderBatteryInfo = (batteryInfo: BatteryInfo) => {
    const { label, value } = batteryInfo;
    return (
      <Stack direction="row" spacing={5}>
        <Typography variant="body1" color="initial">
          {label}:
        </Typography>
        <Typography variant="body1" color="initial">
          {value}
        </Typography>
      </Stack>
    );
  };

  render() {
    const {
      batteryInfo,
      solarHealth,
      nodeHealth,
      switchHealth,
      controllerHealth,
      batteryHealth,
      backhaulHealth,
    } = this.props;

    return (
      <>
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
                <Stack direction="column" spacing={2}>
                  <Typography variant="h6" color="initial">
                    Battery information
                  </Typography>
                  {batteryInfo.map(this.renderBatteryInfo)}
                </Stack>
              </Paper>
              <Paper variant="outlined" sx={{ p: 2 }}>
                <Stack direction="column" spacing={2}>
                  <Typography variant="h6" color="initial">
                    Battery information
                  </Typography>
                  {this.renderBatteryInfo({
                    label: 'Charge',
                    value: `${batteryInfo[2].value} %`,
                  })}
                  {this.renderBatteryInfo({
                    label: 'Voltage',
                    value: `${batteryInfo[0].value} V`,
                  })}
                  {this.renderBatteryInfo({
                    label: 'Current',
                    value: `${batteryInfo[1].value} A`,
                  })}
                  {this.renderBatteryInfo({
                    label: 'Power',
                    value: `${batteryInfo[3].value} W`,
                  })}
                </Stack>
              </Paper>
            </Stack>
          </Grid>
        </Grid>
      </>
    );
  }
}

export default SiteOverallHealth;
