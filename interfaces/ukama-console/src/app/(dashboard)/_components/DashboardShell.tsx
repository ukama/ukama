/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useEffect, useState } from 'react';
import UMark from '@/components/UMark';
import { useUiPrefs } from '@/lib/store';
import CommandPalette from './CommandPalette';
import Sidebar from './Sidebar';
import TopBar from './TopBar';

export default function DashboardShell({
  children,
}: {
  children: React.ReactNode;
}) {
  const rail = useUiPrefs((s) => s.rail);
  const [paletteOpen, setPaletteOpen] = useState(false);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault();
        setPaletteOpen((o) => !o);
      }
    };
    document.addEventListener('keydown', onKey);
    return () => document.removeEventListener('keydown', onKey);
  }, []);

  return (
    <div className="app" data-rail={rail}>
      <div className="brandcell">
        <UMark className="umark" />
        <span className="word">ukama</span>
      </div>
      <TopBar />
      <Sidebar />
      <main className="main">{children}</main>
      <CommandPalette open={paletteOpen} onClose={() => setPaletteOpen(false)} />
    </div>
  );
}
