import React from 'react';
import { Grid, Typography, Stack } from '@mui/material';
import { RoundedCard } from '@/styles/global';

interface SiteDetailsProps {
  dateCreated: string;
  location: string;
  numberOfNodes: number;
}

const SiteDetailsCard: React.FC<SiteDetailsProps> = ({
  dateCreated,
  location,
  numberOfNodes,
}) => {
  return (
    <Grid item xs={3}>
      <RoundedCard>
        <Typography variant="h6" gutterBottom sx={{ py: 1 }}>
          Site details
        </Typography>
        <Stack direction="column" spacing={2}>
          <Stack direction="row" spacing={2} alignItems="center">
            <Typography variant="subtitle1">Date created:</Typography>
            <Typography variant="body2">{dateCreated}</Typography>
          </Stack>
          <Stack direction="row" spacing={2} alignItems="center">
            <Typography variant="subtitle1">Location:</Typography>
            <Typography variant="body2">{location}</Typography>
          </Stack>
          <Stack direction="row" spacing={2} alignItems="center">
            <Typography variant="subtitle1">Nodes:</Typography>
            <Typography variant="body2">{numberOfNodes}</Typography>
          </Stack>
        </Stack>
      </RoundedCard>
    </Grid>
  );
};

export default SiteDetailsCard;
