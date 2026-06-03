/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { usePathname, useRouter } from 'next/navigation';
import { LENSES, lensFromPath } from '../_config/nav';
import { Ic } from './icons';

/** Business · Network · Customer switch (the three lenses, BUILD-PLAN §2). */
export default function LensSegment() {
  const pathname = usePathname();
  const router = useRouter();
  const lens = lensFromPath(pathname);

  return (
    <div className="viewseg" title="Switch console view">
      {LENSES.map((l) => (
        <button
          key={l.id}
          type="button"
          className={lens === l.id ? 'on' : ''}
          onClick={() => router.push(l.href)}
        >
          <Ic name={l.icon} sx={{ fontSize: 17 }} />
          <span className="vlabel">{l.label}</span>
        </button>
      ))}
    </div>
  );
}
