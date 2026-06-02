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
    <Stack spacing={1} sx={{ p: 3 }}>
      <Skeleton variant="rectangular" height={56} sx={{ borderRadius: 1 }} />
      {[...Array(8)].map((_, i) => (
        <Skeleton key={i} variant="rectangular" height={52} sx={{ borderRadius: 1 }} />
      ))}
    </Stack>
  );
}
