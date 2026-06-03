/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** "Today" date-range filter chip with dropdown (biz-common.jsx). */
import { useState } from 'react';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import CalendarTodayRounded from '@mui/icons-material/CalendarTodayRounded';
import CheckRounded from '@mui/icons-material/CheckRounded';
import ExpandMoreRounded from '@mui/icons-material/ExpandMoreRounded';

const DEFAULT_OPTIONS = ['Today', 'Last 7 days', 'This month', 'This quarter'];

export default function DateChip({
  value: controlled,
  options = DEFAULT_OPTIONS,
  onChange,
}: {
  value?: string;
  options?: string[];
  onChange?: (value: string) => void;
}) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const [internal, setInternal] = useState(options[0] ?? 'Today');
  const value = controlled ?? internal;

  return (
    <>
      <button
        type="button"
        className="chip-filter"
        style={{ height: 36 }}
        onClick={(e) => setAnchor(e.currentTarget)}
      >
        <CalendarTodayRounded sx={{ fontSize: 16 }} />
        {value}
        <ExpandMoreRounded sx={{ fontSize: 18, color: 'var(--uk-ink-3)' }} />
      </button>
      <Menu
        anchorEl={anchor}
        open={!!anchor}
        onClose={() => setAnchor(null)}
        slotProps={{ paper: { sx: { width: 170, mt: 0.5 } } }}
      >
        {options.map((o) => (
          <MenuItem
            key={o}
            selected={o === value}
            sx={{ fontSize: 13.5 }}
            onClick={() => {
              setInternal(o);
              onChange?.(o);
              setAnchor(null);
            }}
          >
            {o}
            {o === value && (
              <CheckRounded sx={{ fontSize: 17, color: 'primary.main', ml: 'auto' }} />
            )}
          </MenuItem>
        ))}
      </Menu>
    </>
  );
}
