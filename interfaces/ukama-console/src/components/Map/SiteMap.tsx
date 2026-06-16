/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Site status dot, shared by the Business Home sites list/map. */
import { BIZ_DOT } from './siteMapShared';
import type { SiteMapSite } from './siteMapShared';

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
