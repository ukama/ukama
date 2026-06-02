/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Button, Stack, Typography } from '@mui/material';
import Link from 'next/link';

export default function NotFound() {
  return (
    <Stack
      spacing={3}
      alignItems="center"
      justifyContent="center"
      sx={{ minHeight: '100vh' }}
    >
      <Typography variant="h2" fontWeight={700} color="text.secondary">
        404
      </Typography>
      <Typography variant="h5">Page not found</Typography>
      <Typography variant="body2" color="text.secondary">
        The page you are looking for does not exist.
      </Typography>
      <Button component={Link} href="/console/home" variant="contained">
        Go to dashboard
      </Button>
    </Stack>
  );
}
