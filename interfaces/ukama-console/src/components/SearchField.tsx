/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Search field — MUI TextField themed to the design `.field` look. */
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import SearchRounded from '@mui/icons-material/SearchRounded';

export default function SearchField({
  value,
  onChange,
  placeholder,
  width = 240,
}: {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  width?: number | string;
}) {
  return (
    <TextField
      value={value}
      onChange={(e) => onChange(e.target.value)}
      placeholder={placeholder}
      sx={{ width }}
      slotProps={{
        input: {
          startAdornment: (
            <InputAdornment position="start">
              <SearchRounded sx={{ fontSize: 19, color: 'var(--uk-ink-3)' }} />
            </InputAdornment>
          ),
        },
      }}
    />
  );
}
