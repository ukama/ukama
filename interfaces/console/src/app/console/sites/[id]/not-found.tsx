/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Button, Stack, Typography } from '@mui/material';
import Link from 'next/link';

export default function NotFound() {
  return (
    <Stack spacing={2} alignItems="center" justifyContent="center" sx={{ minHeight: '60vh' }}>
      <Typography variant="h5">Site not found</Typography>
      <Typography variant="body2" color="text.secondary">
        This site may have been removed or the ID is incorrect.
      </Typography>
      <Button component={Link} href="/console/sites" variant="contained">
        Back to sites
      </Button>
    </Stack>
  );
}
