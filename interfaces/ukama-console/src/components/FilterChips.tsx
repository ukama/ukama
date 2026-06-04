/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Filter chip row — clickable MUI Chips with the chipFilter variant. */
import Chip from '@mui/material/Chip';
import Stack from '@mui/material/Stack';

export interface FilterChipOption {
  value: string;
  label: string;
  count?: number;
}

const ACTIVE_SX = {
  borderColor: 'var(--uk-ac)',
  color: 'var(--uk-ac-dark)',
  background: 'var(--uk-ac-soft)',
  '&:hover': { background: 'var(--uk-ac-soft)' },
} as const;

export default function FilterChips({
  options,
  value,
  onChange,
}: {
  options: FilterChipOption[];
  value: string;
  onChange: (value: string) => void;
}) {
  return (
    <Stack direction="row" gap={1} flexWrap="wrap">
      {options.map((o) => (
        <Chip
          key={o.value}
          variant="chipFilter"
          clickable
          onClick={() => onChange(o.value)}
          sx={o.value === value ? ACTIVE_SX : undefined}
          label={
            <>
              {o.label}
              {o.count != null && (
                <span className="tnum" style={{ opacity: 0.6 }}>
                  {o.count}
                </span>
              )}
            </>
          }
        />
      ))}
    </Stack>
  );
}
