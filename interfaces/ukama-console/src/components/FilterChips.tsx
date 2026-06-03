/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Filter chip row — status filters with optional counts (chip-filter /
 * PillFilter in the prototype; one component for both lenses).
 */
export interface FilterChipOption {
  value: string;
  label: string;
  count?: number;
}

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
    <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
      {options.map((o) => (
        <button
          key={o.value}
          type="button"
          className={`chip-filter${o.value === value ? ' on' : ''}`}
          onClick={() => onChange(o.value)}
        >
          {o.label}
          {o.count != null && (
            <span className="tnum" style={{ opacity: 0.6 }}>
              {o.count}
            </span>
          )}
        </button>
      ))}
    </div>
  );
}
