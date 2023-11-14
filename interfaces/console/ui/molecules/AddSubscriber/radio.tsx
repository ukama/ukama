import React from 'react';
import { Paper, Stack, Radio, Typography } from '@mui/material';
import colors from '@/styles/theme/colors';

interface CustomRadioButtonProps {
  label: string;
  count: number | undefined;
  selected: boolean;
  onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

const CustomRadioButton: React.FC<CustomRadioButtonProps> = ({
  label,
  count,
  selected,
  onChange,
}) => {
  return (
    <Paper variant="outlined" sx={{}}>
      <Stack direction="row" spacing={1} alignItems="center">
        <Radio
          checked={selected}
          onChange={onChange}
          inputProps={{
            'aria-label': label,
          }}
        />
        <Stack direction="row" spacing={1} alignItems="center">
          <Typography variant="body1" sx={{ color: colors.black }}>
            {label}
          </Typography>
          <Typography variant="body2" sx={{ color: colors.black70 }}>
            {` (${count || 0} left)`}
          </Typography>
        </Stack>
      </Stack>
    </Paper>
  );
};

export default CustomRadioButton;
