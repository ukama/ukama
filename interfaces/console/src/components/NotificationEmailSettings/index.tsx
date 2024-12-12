/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Paper, Typography, Stack, TextField, Box } from '@mui/material';
import colors from '@/theme/colors';

interface NotificationEmailProps {
  primaryEmail: string;
  additionalEmails?: string[];
}

const NotificationEmailSettings: React.FC<NotificationEmailProps> = ({
  primaryEmail,
  additionalEmails = [],
}) => {
  return (
    <Paper
      elevation={2}
      sx={{
        p: 4,
        mt: 2,
        borderRadius: '10px',
        bgcolor: colors.white,
      }}
    >
      <Typography variant="h6">Notification Settings</Typography>
      <Stack direction="column" spacing={2} sx={{ mt: 2 }}>
        <Typography variant="body2" sx={{ color: colors.black54 }}>
          All entered emails will receive receipts for the monthly bill.{' '}
        </Typography>

        <TextField
          fullWidth
          label="PRIMARY EMAIL"
          value={primaryEmail}
          variant="outlined"
          disabled
          InputLabelProps={{ shrink: true }}
        />

        {additionalEmails.map((email) => (
          <Box
            key={email}
            sx={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              bgcolor: colors.black38,
              p: 1,
              borderRadius: 1,
            }}
          >
            <Typography variant="body2">{email}</Typography>
          </Box>
        ))}
      </Stack>
    </Paper>
  );
};

export default NotificationEmailSettings;
