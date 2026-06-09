/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Entity switcher under detail titles (node-site-detail.jsx DetailPicker). */
import { useState } from 'react';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import UnfoldMoreRounded from '@mui/icons-material/UnfoldMoreRounded';
import StatusBadge from '@/components/StatusBadge';

export interface PickerItem {
  id: string;
  label: string;
  status: string;
}

export default function DetailPicker({
  value,
  items,
  onPick,
}: {
  value: PickerItem;
  items: PickerItem[];
  onPick: (item: PickerItem) => void;
}) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);

  return (
    <>
      <button type="button" className="detail-picker" onClick={(e) => setAnchor(e.currentTarget)}>
        <span className="tnum">{value.label}</span>
        <UnfoldMoreRounded sx={{ fontSize: 18, color: 'var(--uk-ink-3)' }} />
      </button>
      <Menu
        anchorEl={anchor}
        open={!!anchor}
        onClose={() => setAnchor(null)}
        slotProps={{ paper: { sx: { width: 250, mt: 0.5 } } }}
      >
        {items.map((it) => (
          <MenuItem
            key={it.id}
            selected={it.id === value.id}
            sx={{ fontSize: 13, justifyContent: 'space-between', gap: 1.5 }}
            onClick={() => {
              onPick(it);
              setAnchor(null);
            }}
          >
            <span className="tnum">{it.label}</span>
            <StatusBadge status={it.status} />
          </MenuItem>
        ))}
      </Menu>
    </>
  );
}
