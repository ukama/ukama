import React from 'react';
import { Paper, Radio, Typography, Stack, Grid } from '@mui/material';
import { colors } from '@/styles/theme';

interface SimTypeProps {
  simType: string;
  label: string;
  count: number | undefined;
  selectedSimType: string;
  handleSimTypeChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

const SimTypeComponent: React.FC<SimTypeProps> = ({
  simType,
  label,
  count,
  selectedSimType,
  handleSimTypeChange,
}) => (
  <Grid item xs={6}>
    <Paper variant="outlined" sx={{}}>
      <Stack direction="row" spacing={1} alignItems="center">
        <Radio
          value={simType}
          name={simType}
          onChange={handleSimTypeChange}
          checked={selectedSimType === simType}
          inputProps={{
            'aria-label': simType,
          }}
        />
        <Stack direction="row" spacing={1} alignItems="center">
          <Typography variant="body1" sx={{ color: colors.black }}>
            {`${label} `}
          </Typography>
          <Typography variant="body2" sx={{ color: colors.black70 }}>
            {` (${count || 0} left)`}
          </Typography>
        </Stack>
      </Stack>
    </Paper>
  </Grid>
);

export default SimTypeComponent;
