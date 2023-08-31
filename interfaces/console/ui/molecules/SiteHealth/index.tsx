import React from 'react';
import { Stack, Typography, Paper, Grid } from '@mui/material';
import { SiteHealth } from '@/public/svg';

interface SiteOverallHealthProps {
  voltage: string;
  current: string;
  power: string;
  modelNumber: string;
  version: string;
  charge: string;
  solarHealth: 'good' | 'warning';
  nodeHealth: 'good' | 'warning';
  switchHealth: 'good' | 'warning';
  controllerHealth: 'good' | 'warning';
  batteryHealth: 'good' | 'warning';
  backhaulHealth: 'good' | 'warning';
}

const SiteOverallHealth: React.FC<SiteOverallHealthProps> = ({
  voltage,
  current,
  power,
  modelNumber,
  version,
  charge,
  solarHealth,
  nodeHealth,
  switchHealth,
  controllerHealth,
  batteryHealth,
  backhaulHealth,
}) => {
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
                <Stack direction="row" spacing={5}>
                  <Typography variant="body1" color="initial">
                    Model number :
                  </Typography>
                  <Typography variant="body1" color="initial">
                    {modelNumber}
                  </Typography>
                </Stack>
                <Stack direction="row" spacing={5}>
                  <Typography variant="body1" color="initial">
                    Version :
                  </Typography>
                  <Typography variant="body1" color="initial">
                    {` V ${version}`}
                  </Typography>
                </Stack>
              </Stack>
            </Paper>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Stack direction="column" spacing={2}>
                <Typography variant="h6" color="initial">
                  Battery information
                </Typography>
                <Stack direction="row" spacing={5}>
                  <Typography variant="body1" color="initial">
                    Charge :
                  </Typography>
                  <Typography variant="body1" color="initial">
                    {`${charge} %`}
                  </Typography>
                </Stack>
                <Stack direction="row" spacing={5}>
                  <Typography variant="body1" color="initial">
                    Voltage :
                  </Typography>
                  <Typography variant="body1" color="initial">
                    {`${voltage} V`}
                  </Typography>
                </Stack>
                <Stack direction="row" spacing={5}>
                  <Typography variant="body1" color="initial">
                    Current :
                  </Typography>
                  <Typography variant="body1" color="initial">
                    {`${current} A`}
                  </Typography>
                </Stack>
                <Stack direction="row" spacing={5}>
                  <Typography variant="body1" color="initial">
                    Power :
                  </Typography>
                  <Typography variant="body1" color="initial">
                    {`${power} W`}
                  </Typography>
                </Stack>
              </Stack>
            </Paper>
          </Stack>
        </Grid>
      </Grid>
    </>
  );
};

export default SiteOverallHealth;
