/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useEffect } from 'react';
import { Button, Stack, Typography } from '@mui/material';

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function Error({ error, reset }: ErrorProps) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <Stack
      spacing={2}
      alignItems="center"
      justifyContent="center"
      sx={{ minHeight: '60vh', p: 3 }}
    >
      <Typography variant="h6" color="error">
        Something went wrong
      </Typography>
      <Typography variant="body2" color="text.secondary" textAlign="center">
        {error.message ?? 'An unexpected error occurred'}
      </Typography>
      <Button variant="contained" onClick={reset}>
        Try again
      </Button>
    </Stack>
  );
}
