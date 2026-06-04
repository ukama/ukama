/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Route-level skeleton while a dashboard segment loads (BUILD-PLAN §13.2). */
import Skeleton from '@mui/material/Skeleton';

export default function DashboardLoading() {
  return (
    <div className="page">
      <Skeleton variant="text" width={220} height={42} />
      <Skeleton variant="text" width={340} height={20} sx={{ mb: 2 }} />
      <Skeleton variant="rounded" height={108} sx={{ mb: 2 }} />
      <Skeleton variant="rounded" height={320} />
    </div>
  );
}
