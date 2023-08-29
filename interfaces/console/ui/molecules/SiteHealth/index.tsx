import React from 'react';
import { Stack, Typography, Divider, Box, Grid } from '@mui/material';

import {
  NodeIcon,
  SvgContainer,
  SwitchIcon,
  ControllerIcon,
} from '@/public/svg';

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
        <Grid item xs={6} container spacing={2}>
          <Grid item xs={2}>
            <SvgContainer>
              <NodeIcon />
            </SvgContainer>
          </Grid>
          <Grid item xs={2}>
            <SvgContainer>
              <SwitchIcon />
            </SvgContainer>
          </Grid>
          <Grid item xs={2}>
            <SvgContainer>
              <ControllerIcon />
            </SvgContainer>
          </Grid>
        </Grid>
        <Grid item xs={6}></Grid>
      </Grid>
    </>
  );
};

export default BaseStationSiteHealth;
