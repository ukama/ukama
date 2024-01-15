import React from 'react';
import { Grid, Typography, Box, Stack } from '@mui/material';
import { RoundedCard } from '@/styles/global';
import { PersonIcon, TowerIcon, NodeIcon, BatteryIcon } from '../SvgIcons';
import ErrorIcon from '@mui/icons-material/Error';
import { colors } from '@/styles/theme';

interface Site {
  name: string;
  details: string;
  batteryStatus: 'charging' | 'notCharging';
  nodeStatus: 'online' | 'offline';
  towerStatus: 'online' | 'offline';
  numberOfPersonsConnected: number;
}

interface SiteCardProps {
  sites: Site[];
}

const SiteCard: React.FC<SiteCardProps> = ({ sites }) => {
  return (
    <RoundedCard>
      {sites.map((site, index) => (
        <Grid container spacing={1} key={index} alignItems={'center'}>
          <Grid item xs={12}>
            <Stack direction="row" spacing={1} alignItems={'center'}>
              <Typography variant="h5">{site.name}</Typography>
              {(site.towerStatus === 'offline' ||
                site.nodeStatus === 'offline' ||
                site.batteryStatus === 'notCharging') && (
                <ErrorIcon sx={{ color: colors.red, fontSize: 18 }} />
              )}
            </Stack>

            <Typography variant="body2">{site.details}</Typography>
          </Grid>
          <Grid item xs={2}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <PersonIcon />
              <Typography variant="caption">
                {`${site.numberOfPersonsConnected}`}
              </Typography>
            </Stack>
          </Grid>
          <Grid item xs={3}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <NodeIcon status={site.nodeStatus} />
              <Typography variant="caption">{`${site.nodeStatus}`}</Typography>
            </Stack>
          </Grid>
          <Grid item xs={4}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <BatteryIcon status={site.batteryStatus} />
              <Typography variant="caption">{`${site.batteryStatus}`}</Typography>
            </Stack>
          </Grid>
          <Grid item xs={3}>
            <Stack direction={'row'} spacing={1} alignItems={'center'}>
              <TowerIcon status={site.towerStatus} />
              <Typography variant="caption">{`${site.towerStatus}`}</Typography>
            </Stack>
          </Grid>
        </Grid>
      ))}
    </RoundedCard>
  );
};

export default SiteCard;
