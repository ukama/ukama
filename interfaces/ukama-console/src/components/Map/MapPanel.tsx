/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Lazy-loaded network map (BUILD-PLAN §7.1/§14 — maps are code-split). */
import dynamic from 'next/dynamic';
import Skeleton from '@mui/material/Skeleton';

const MapPanel = dynamic(() => import('./MapPanelImpl'), {
  ssr: false,
  loading: () => (
    <Skeleton variant="rounded" sx={{ width: '100%', height: '100%', minHeight: 200 }} />
  ),
});

export default MapPanel;
