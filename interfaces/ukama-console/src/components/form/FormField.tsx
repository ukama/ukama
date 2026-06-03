/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Form primitives matching the prototype's dialog fields (form-dialogs.jsx),
 * designed for react-hook-form `register` spreading (BUILD-PLAN §3 forms).
 */
import { forwardRef } from 'react';
import ExpandMoreRounded from '@mui/icons-material/ExpandMoreRounded';

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
  React.InputHTMLAttributes<HTMLInputElement> & {
    prefix?: string;
    invalid?: boolean;
  }
>(function TextInput({ prefix, invalid, ...props }, ref) {
  return (
    <div className={`field${invalid ? ' invalid' : ''}`}>
      {prefix && <span style={{ color: 'var(--uk-ink-3)', fontSize: 14 }}>{prefix}</span>}
      <input ref={ref} {...props} />
    </div>
  );
});

export const SelectInput = forwardRef<
  HTMLSelectElement,
  React.SelectHTMLAttributes<HTMLSelectElement> & {
    invalid?: boolean;
    placeholder?: string;
    options: { value: string; label: string }[];
  }
>(function SelectInput({ invalid, placeholder, options, ...props }, ref) {
  return (
    <div className={`field${invalid ? ' invalid' : ''}`}>
      <select ref={ref} {...props}>
        {placeholder && <option value="">{placeholder}</option>}
        {options.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
      <ExpandMoreRounded sx={{ fontSize: 19, color: 'var(--uk-ink-3)', pointerEvents: 'none' }} />
    </div>
  );
});
