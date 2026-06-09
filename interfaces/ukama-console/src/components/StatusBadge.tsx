/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * One status vocabulary, everywhere (design finding #4) — MUI Chip with
 * custom tone variants from the theme (§7.2 C).
 * `variant="dot"` = ops badge with leading dot; `variant="pill"` = no dot.
 */
import Chip from '@mui/material/Chip';
import { STATUS_MAP } from '@/data';
import type { StatusTone } from '@/data';

type ToneVariant = 'ok' | 'warn' | 'err' | 'neut' | 'info';

const TONE_VARIANT: Record<StatusTone, ToneVariant> = {
  ok: 'ok',
  warn: 'warn',
  err: 'err',
  neutral: 'neut',
};

const DOT_COLOR: Record<ToneVariant, string> = {
  ok: 'var(--uk-success-bright)',
  warn: 'var(--uk-warning)',
  err: 'var(--uk-error)',
  neut: 'var(--uk-ink-3)',
  info: 'var(--uk-ac)',
};

/** Extra statuses used by the business pills (biz-common.jsx BIZ_PILL). */
const EXTRA: Record<string, { tone: ToneVariant; label: string }> = {
  deployed: { tone: 'ok', label: 'Deployed' },
  testing: { tone: 'warn', label: 'Testing' },
  expired: { tone: 'warn', label: 'Expired' },
  assigned: { tone: 'info', label: 'Assigned' },
  available: { tone: 'info', label: 'Available' },
  lowsales: { tone: 'neut', label: 'Low sales' },
  warning: { tone: 'warn', label: 'Warning' },
  paid: { tone: 'ok', label: 'Paid' },
  rma: { tone: 'err', label: 'RMA' },
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
  const tone: ToneVariant = meta ? TONE_VARIANT[meta.tone] : (extra?.tone ?? 'neut');
  const label = children ?? (meta ? meta.label : (extra?.label ?? status));

  return (
    <Chip
      variant={tone}
      label={label}
      {...(variant === 'dot'
        ? {
            icon: (
              <span
                style={{
                  width: 7,
                  height: 7,
                  borderRadius: '50%',
                  background: DOT_COLOR[tone],
                  display: 'inline-block',
                }}
              />
            ),
          }
        : {})}
    />
  );
}
