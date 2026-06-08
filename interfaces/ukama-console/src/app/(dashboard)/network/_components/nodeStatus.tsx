/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Shared node connectivity dot + state chip (cards + detail page). */

/** Connectivity dot color from the raw value. */
export const connColor = (c?: string): string => {
  const v = (c ?? '').toLowerCase();
  if (v === 'online') return 'var(--uk-success-bright)';
  if (v === 'offline') return 'var(--uk-error)';
  return 'var(--uk-ink-3)';
};

/**
 * Node state chip styling: Unknown = plain (neutral), Configured = green,
 * Operational = blue, Faulty = matt red.
 */
const STATE_STYLE: Record<
  string,
  { bg: string; color: string; border: string; label: string }
> = {
  unknown: {
    bg: 'transparent',
    color: 'var(--uk-ink)',
    border: 'var(--uk-line)',
    label: 'Unknown',
  },
  configured: {
    bg: 'rgba(29, 205, 159, 0.14)',
    color: 'var(--uk-success-bright)',
    border: 'transparent',
    label: 'Configured',
  },
  operational: {
    bg: 'var(--uk-ac-soft)',
    color: 'var(--uk-ac-dark)',
    border: 'transparent',
    label: 'Operational',
  },
  faulty: {
    bg: 'rgba(207, 18, 27, 0.14)',
    color: '#e2575f',
    border: 'transparent',
    label: 'Faulty',
  },
};

const stateStyle = (s?: string) => {
  const v = (s ?? '').toLowerCase();
  return (
    STATE_STYLE[v] ?? {
      ...STATE_STYLE.unknown!,
      label: v ? v[0]!.toUpperCase() + v.slice(1) : 'Unknown',
    }
  );
};

export function ConnectivityDot({ connectivity }: { connectivity?: string }) {
  return (
    <span
      title={`Connectivity: ${connectivity ?? 'Unknown'}`}
      style={{
        width: 8,
        height: 8,
        borderRadius: '50%',
        background: connColor(connectivity),
        flex: 'none',
        display: 'inline-block',
      }}
    />
  );
}

export function StateChip({ state }: { state?: string }) {
  const s = stateStyle(state);
  return (
    <span
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        fontSize: 12,
        fontWeight: 600,
        padding: '3px 10px',
        borderRadius: 999,
        background: s.bg,
        color: s.color,
        border: `1px solid ${s.border}`,
        whiteSpace: 'nowrap',
      }}
    >
      {s.label}
    </span>
  );
}
