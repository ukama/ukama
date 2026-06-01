/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Skeleton, Stack } from '@mui/material';

export default function Loading() {
  return (
    <Stack spacing={2} sx={{ p: 3 }}>
      <Skeleton variant="rectangular" height={48} width={300} />
      <Stack direction="row" spacing={2}>
        <Skeleton variant="rectangular" width="40%" height={300} sx={{ borderRadius: 2 }} />
        <Skeleton variant="rectangular" width="60%" height={300} sx={{ borderRadius: 2 }} />
      </Stack>
    </Stack>
  );
}
