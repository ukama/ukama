/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Skeleton } from '@mui/material';
import dynamic from 'next/dynamic';

const LineChart = dynamic(() => import('./index'), {
  ssr: false,
  loading: () => (
    <Skeleton variant="rectangular" width="100%" height={200} sx={{ borderRadius: 2 }} />
  ),
});

export default LineChart;
