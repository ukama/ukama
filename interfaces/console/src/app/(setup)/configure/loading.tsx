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
    <Stack direction="row" spacing={2} sx={{ p: 3, height: '100vh' }}>
      <Skeleton variant="rectangular" width="35%" sx={{ borderRadius: 2 }} />
      <Skeleton variant="rectangular" width="65%" sx={{ borderRadius: 2 }} />
    </Stack>
  );
}
