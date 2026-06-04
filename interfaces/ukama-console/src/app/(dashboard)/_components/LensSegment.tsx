/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import CheckRounded from '@mui/icons-material/CheckRounded';
import ExpandMoreRounded from '@mui/icons-material/ExpandMoreRounded';
import { LENSES, lensFromPath } from '../_config/nav';
import { Ic } from './icons';

/** Business · Network · Customer switch (the three lenses, BUILD-PLAN §2).
 *  Desktop top-bar segmented control; hidden ≤900px (see LensDropdown). */
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

/** Mobile lens switch — dropdown at the top of the nav drawer. */
export function LensDropdown({ onNavigate }: { onNavigate?: () => void }) {
  const pathname = usePathname();
  const router = useRouter();
  const lens = lensFromPath(pathname);
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const current = LENSES.find((l) => l.id === lens) ?? LENSES[0];

  return (
    <>
      <button
        type="button"
        className="lens-mobile"
        aria-label="Switch console view"
        onClick={(e) => setAnchor(e.currentTarget)}
      >
        {current && <Ic name={current.icon} sx={{ fontSize: 18, color: 'var(--uk-ac)' }} />}
        {current?.label}
        <ExpandMoreRounded sx={{ fontSize: 18, marginLeft: 'auto', opacity: 0.6 }} />
      </button>
      <Menu
        anchorEl={anchor}
        open={!!anchor}
        onClose={() => setAnchor(null)}
        slotProps={{ paper: { sx: { width: 232, mt: 0.5 } } }}
      >
        {LENSES.map((l) => (
          <MenuItem
            key={l.id}
            selected={l.id === lens}
            sx={{ fontSize: 13.5, gap: 1.2 }}
            onClick={() => {
              setAnchor(null);
              onNavigate?.();
              router.push(l.href);
            }}
          >
            <Ic name={l.icon} sx={{ fontSize: 18, color: 'var(--uk-ink-3)' }} />
            {l.label}
            {l.id === lens && (
              <CheckRounded sx={{ fontSize: 17, color: 'primary.main', ml: 'auto' }} />
            )}
          </MenuItem>
        ))}
      </Menu>
    </>
  );
}
