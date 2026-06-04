/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Form primitives — MUI TextField/Select themed to the design dialog
 * fields, ready for react-hook-form `register` spreading (BUILD-PLAN §3).
 */
import { forwardRef } from 'react';
import InputAdornment from '@mui/material/InputAdornment';
import type { InputBaseComponentProps } from '@mui/material/InputBase';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import type { TextFieldProps } from '@mui/material/TextField';

export function Field({
  label,
  required,
  hint,
  error,
  children,
}: {
  label: string;
  required?: boolean;
  hint?: string;
  error?: string;
  children: React.ReactNode;
}) {
  return (
    <label className="ff">
      <div className="ff-label">
        {label}
        {required && <span style={{ color: 'var(--uk-error)' }}> *</span>}
      </div>
      {children}
      {error ? (
        <div className="ff-err">{error}</div>
      ) : hint ? (
        <div className="ff-hint">{hint}</div>
      ) : null}
    </label>
  );
}

export const TextInput = forwardRef<
  HTMLInputElement,
  Omit<TextFieldProps, 'variant'> & { prefix?: string; invalid?: boolean }
>(function TextInput({ prefix, invalid, ...props }, ref) {
  return (
    <TextField
      inputRef={ref}
      fullWidth
      error={invalid}
      slotProps={{
        input: {
          startAdornment: prefix ? (
            <InputAdornment position="start">
              <span style={{ color: 'var(--uk-ink-3)', fontSize: 14 }}>{prefix}</span>
            </InputAdornment>
          ) : undefined,
        },
      }}
      {...props}
    />
  );
});

export const SelectInput = forwardRef<
  HTMLSelectElement,
  React.SelectHTMLAttributes<HTMLSelectElement> & {
    invalid?: boolean;
    placeholder?: string;
    options: { value: string; label: string }[];
  }
>(function SelectInput({ invalid, placeholder, options, className: _c, style: _s, ...props }, ref) {
  return (
    <TextField
      fullWidth
      select
      error={invalid}
      defaultValue=""
      slotProps={{
        select: {
          native: true,
          inputRef: ref,
          inputProps: props as InputBaseComponentProps,
        },
      }}
    >
      {placeholder && <option value="">{placeholder}</option>}
      {options.map((o) => (
        <option key={o.value} value={o.value}>
          {o.label}
        </option>
      ))}
    </TextField>
  );
});

export { MenuItem };
