/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Skeleton } from '@mui/material';
import Grid from '@mui/material/Grid';

export default function Loading() {
  return (
    <Grid container spacing={2} sx={{ p: 3 }}>
      {[...Array(3)].map((_, i) => (
        <Grid key={i} size={{ xs: 12, md: 4 }}>
          <Skeleton variant="rectangular" height={140} sx={{ borderRadius: 2 }} />
        </Grid>
      ))}
      <Grid size={12}>
        <Skeleton variant="rectangular" height={300} sx={{ borderRadius: 2 }} />
      </Grid>
    </Grid>
  );
}
