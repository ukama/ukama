/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * One status vocabulary, everywhere (design finding #4).
 * `variant="dot"` = ops badge with leading dot (prototype Badge).
 * `variant="pill"` = business pill without dot (prototype BizPill).
 */
import { STATUS_MAP } from '@/data';
import type { StatusTone } from '@/data';

const TONE_CLS: Record<StatusTone | 'info', string> = {
  ok: 'badge-ok',
  warn: 'badge-warn',
  err: 'badge-err',
  neutral: 'badge-neut',
  info: 'badge-info',
};

/** Extra statuses used by the business pills (biz-common.jsx BIZ_PILL). */
const EXTRA: Record<string, { cls: string; label: string }> = {
  deployed: { cls: 'badge-ok', label: 'Deployed' },
  testing: { cls: 'badge-warn', label: 'Testing' },
  expired: { cls: 'badge-warn', label: 'Expired' },
  assigned: { cls: 'badge-info', label: 'Assigned' },
  available: { cls: 'badge-info', label: 'Available' },
  lowsales: { cls: 'badge-neut', label: 'Low sales' },
  warning: { cls: 'badge-warn', label: 'Warning' },
  paid: { cls: 'badge-ok', label: 'Paid' },
  rma: { cls: 'badge-err', label: 'RMA' },
};

export default function StatusBadge({
  status,
  variant = 'dot',
  children,
}: {
  status: string;
  variant?: 'dot' | 'pill';
  children?: React.ReactNode;
}) {
  const meta = STATUS_MAP[status];
  const extra = EXTRA[status];
  const cls = meta ? TONE_CLS[meta.tone] : extra ? extra.cls : 'badge-neut';
  const label = children ?? (meta ? meta.label : extra ? extra.label : status);

  return (
    <span
      className={`badge ${cls}`}
      style={variant === 'pill' ? { paddingLeft: 9 } : undefined}
    >
      {variant === 'dot' && <span className="dot" />}
      {label}
    </span>
  );
}
