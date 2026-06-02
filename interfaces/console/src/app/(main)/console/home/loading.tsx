/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Skeleton, Stack } from '@mui/material';
import Grid from '@mui/material/Grid2';

export default function Loading() {
  return (
    <Grid container rowSpacing={2} columnSpacing={2}>
      <Grid size={12}>
        <Stack direction="row" spacing={1}>
          <Skeleton variant="rectangular" width="33%" height={90} sx={{ borderRadius: '10px' }} />
          <Skeleton variant="rectangular" width="33%" height={90} sx={{ borderRadius: '10px' }} />
          <Skeleton variant="rectangular" width="33%" height={90} sx={{ borderRadius: '10px' }} />
        </Stack>
      </Grid>
      <Grid size={12}>
        <Skeleton
          variant="rectangular"
          width="100%"
          sx={{ borderRadius: '10px', height: { xs: 'calc(100vh - 280px)', md: 'calc(100vh - 318px)' } }}
        />
      </Grid>
    </Grid>
  );
}
