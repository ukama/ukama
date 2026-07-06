/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useEffect, useState } from 'react';

/** Short human label for an operation type; falls back to the raw type. */
const OP_LABELS: Record<string, string> = {
  restartnode: 'Restart',
  restartsite: 'Site restart',
  togglerf: 'RF change',
  togglerfstatus: 'RF change',
  toggleservice: 'Service change',
};

export function prettyOpType(type?: string | null): string {
  if (!type) return 'Operation';
  return OP_LABELS[type.toLowerCase()] ?? type;
}

/** e.g. "45s", "1m 20s", "2h 05m". */
export function formatElapsed(fromIso?: string | null, now: number = Date.now()): string {
  if (!fromIso) return '';
  const start = Date.parse(fromIso);
  if (!Number.isFinite(start)) return '';
  const s = Math.max(0, Math.floor((now - start) / 1000));
  if (s < 60) return `${s}s`;
  const m = Math.floor(s / 60);
  if (m < 60) return `${m}m ${String(s % 60).padStart(2, '0')}s`;
  const h = Math.floor(m / 60);
  return `${h}h ${String(m % 60).padStart(2, '0')}m`;
}

/**
 * Live-updating elapsed string for an in-progress operation. Ticks every
 * second while `fromIso` is set, and stops (returns '') when it clears.
 */
export function useElapsed(fromIso?: string | null): string {
  const [, tick] = useState(0);
  useEffect(() => {
    if (!fromIso) return;
    const id = setInterval(() => tick((n) => n + 1), 1000);
    return () => clearInterval(id);
  }, [fromIso]);
  return formatElapsed(fromIso);
}

/**
 * Attribution suffix for an operation started by someone else, e.g.
 * " by alan". Empty when unknown.
 */
export function byWhom(requestedBy?: string | null): string {
  return requestedBy ? ` by ${requestedBy}` : '';
}
