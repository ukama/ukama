/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import type { PaymentFragment } from '@/client/graphql/billing.generated';
import { parseTimestamp } from '@/lib/parsers';

/** Decode a payment's base64 metadata (e.g. "eyJzaW0iOiAi..."} ) and return its sim id. */
export function paymentSim(metadata?: string | null): string | null {
  if (!metadata) return null;
  try {
    const decoded =
      typeof atob === 'function'
        ? atob(metadata)
        : Buffer.from(metadata, 'base64').toString('utf-8');
    const json = JSON.parse(decoded) as { sim?: unknown };
    return typeof json.sim === 'string' ? json.sim : null;
  } catch {
    return null;
  }
}

/**
 * Pick the payment that belongs to a specific SIM. `item_id` (the data-plan id)
 * can return payments across SIMs and repeat purchases, so we match exactly on
 * the decoded metadata.sim and take the most recent. Returns undefined when no
 * payment is tied to this SIM (e.g. a directly-allocated package).
 */
export function pickPaymentForSim(
  payments: readonly PaymentFragment[],
  simId: string,
): PaymentFragment | undefined {
  return [...payments]
    .filter((p) => paymentSim(p.metadata) === simId)
    .sort((a, b) => parseTimestamp(b.createdAt) - parseTimestamp(a.createdAt))[0];
}

const cap = (s: string): string =>
  s ? s.charAt(0).toUpperCase() + s.slice(1) : s;

/** Format a payment timestamp as "16 Jul 2026, 17:17 UTC". */
export function formatReceiptDateTime(s?: string | null): string {
  const t = parseTimestamp(s ?? undefined);
  if (Number.isNaN(t)) return '—';
  const fmt = new Intl.DateTimeFormat('en-GB', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    timeZone: 'UTC',
    hour12: false,
  }).format(new Date(t));
  return `${fmt} UTC`;
}

/** Flat, presentation-ready receipt used by both the dialog and the PDF. */
export type ReceiptView = {
  id: string;
  orgName: string;
  planName: string;
  amountLabel: string;
  currency: string;
  method: string;
  status: string;
  paidOn: string;
  simId: string | null;
  description: string;
};

export function buildReceiptView(
  p: PaymentFragment,
  opts: { planName: string; symbol: string; orgName: string },
): ReceiptView {
  const amount = Number(p.amount);
  return {
    id: p.id,
    orgName: opts.orgName,
    planName: opts.planName,
    amountLabel: `${opts.symbol}${Number.isNaN(amount) ? p.amount : amount.toFixed(2)}`,
    currency: (p.currency || '').toUpperCase(),
    method: cap(p.paymentMethod || ''),
    status: cap(p.status || ''),
    paidOn: formatReceiptDateTime(p.paidAt || p.createdAt),
    simId: paymentSim(p.metadata),
    description: p.description || '',
  };
}
