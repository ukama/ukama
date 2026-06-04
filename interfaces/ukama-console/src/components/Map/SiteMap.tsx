/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Lazy-loaded business sites map (BUILD-PLAN §7.1/§14 — code-split). */
import dynamic from 'next/dynamic';
import Skeleton from '@mui/material/Skeleton';
import { BIZ_DOT } from './siteMapShared';
import type { SiteMapSite } from './siteMapShared';

const SiteMap = dynamic(() => import('./SiteMapImpl'), {
  ssr: false,
  loading: () => (
    <div className="card" style={{ padding: 0, overflow: 'hidden', height: '100%', minHeight: 230 }}>
      <Skeleton variant="rounded" sx={{ width: '100%', height: '100%', minHeight: 230 }} />
    </div>
  ),
});

export default SiteMap;
export { BIZ_DOT };
export type { SiteMapSite };

/** Small leading status dot (biz-common.jsx StatusDot). */
export function StatusDot({ status }: { status: SiteMapSite['status'] }) {
  return (
    <span
      style={{
        width: 9,
        height: 9,
        borderRadius: '50%',
        flex: 'none',
        background: BIZ_DOT[status] ?? 'var(--uk-ink-3)',
        display: 'inline-block',
      }}
    />
  );
}
