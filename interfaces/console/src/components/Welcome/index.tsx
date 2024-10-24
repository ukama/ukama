/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import { Button, Paper, Stack, Typography } from '@mui/material';

interface IWelcome {
  orgName: string;
  operatingCountry: string;
  handleNext: () => void;
  handleBack: () => void;
}

export const LField = ({ label, value }: { label: string; value: string }) => {
  return (
    <Stack direction="column" spacing={0.5}>
      <Typography
        variant="caption"
        textTransform={'uppercase'}
        color={colors.black54}
      >
        {label}
      </Typography>
      <Typography variant="body1" fontWeight={400}>
        {value}
      </Typography>
    </Stack>
  );
};

const Welcome = ({
  orgName,
  operatingCountry,
  handleBack,
  handleNext,
}: IWelcome) => {
  return (
    <Paper elevation={0} sx={{ px: 4, py: 2 }}>
      <Typography variant="h6" fontWeight={500}>
        Welcome to Ukama!
      </Typography>
      <Stack direction={'column'} mt={3} mb={3} spacing={2}>
        <Typography variant="body1">
          Welcome to Ukama! Please check to make sure the following details are
          correct before moving onto the next step.
        </Typography>

        <LField label="Organization Name" value={orgName} />
        <LField label="Network Operating Country" value={operatingCountry} />
      </Stack>
      <Stack direction={'row'} justifyContent={'space-between'} spacing={2}>
        <Button
          variant="text"
          onClick={handleBack}
          sx={{ color: colors.black70, p: 0 }}
        >
          Back to Singup
        </Button>
        <Button variant="contained" color="primary" onClick={handleNext}>
          Next
        </Button>
      </Stack>
    </Paper>
  );
};

export default Welcome;
