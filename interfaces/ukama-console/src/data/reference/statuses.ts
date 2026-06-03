/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Canonical status → tone/label mapping (single source of truth for the
 * StatusBadge; ported from the prototype's STATUS_MAP in ds.jsx).
 */

export type StatusTone = 'ok' | 'warn' | 'err' | 'neutral';

export interface StatusMeta {
  tone: StatusTone;
  label: string;
}

export const STATUS_MAP: Record<string, StatusMeta> = {
  online: { tone: 'ok', label: 'Online' },
  healthy: { tone: 'ok', label: 'Healthy' },
  active: { tone: 'ok', label: 'Active' },
  charged: { tone: 'ok', label: 'Charged' },
  degraded: { tone: 'warn', label: 'Degraded' },
  configuring: { tone: 'warn', label: 'Configuring' },
  low: { tone: 'warn', label: 'Low' },
  pending: { tone: 'warn', label: 'Pending' },
  offline: { tone: 'err', label: 'Offline' },
  faulty: { tone: 'err', label: 'Faulty' },
  critical: { tone: 'err', label: 'Critical' },
  failed: { tone: 'err', label: 'Failed' },
  inactive: { tone: 'neutral', label: 'Inactive' },
  idle: { tone: 'neutral', label: 'Idle' },
};
