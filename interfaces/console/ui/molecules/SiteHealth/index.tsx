import React from 'react';
import { Stack, Typography, Divider, Paper, Box, Grid } from '@mui/material';
import Xarrow from 'react-xarrows';
import { SteppedLineTo } from 'react-lineto';

import {
  NodeIcon,
  SvgContainer,
  SwitchIcon,
  ControllerIcon,
  BackHaulIcon,
  BatteryLevelIcon,
  SolarIcon,
  SiteHealth,
} from '@/public/svg';
import { colors } from '@/styles/theme';

interface BaseStationSiteHealthProps {
  batteryLevel: number;
  internetSwitch: boolean;
  controllerSwitch: boolean;
}

const BaseStationSiteHealth: React.FC<BaseStationSiteHealthProps> = ({
  batteryLevel,
  internetSwitch,
  controllerSwitch,
}) => {
  return (
    <>
      <Grid container spacing={2}>
        <Grid item xs={8}>
          <SiteHealth solarHealth={'warning'} />
        </Grid>
        <Grid item xs={4}>
          <Paper>
            <Typography
              variant="body1"
              sx={{ fontWeight: 'semi-bold' }}
              color="initial"
            >
              Battery information
            </Typography>
          </Paper>
        </Grid>
      </Grid>
    </>
  );
};

export default BaseStationSiteHealth;
