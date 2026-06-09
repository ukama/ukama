/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Billing — ported from the design prototype (`data.jsx`). */
import { series } from '@/lib/series';

export interface Invoice {
  id: string;
  period: string;
  amt: number;
  status: 'paid';
}

export const BILLING = {
  current: 1284.5,
  due: '01 Jul 2026',
  cycle: 'Jun 1 – Jun 30, 2026',
  plan: 'Operator · usage based',
  method: 'Visa ·· 4421',
  breakdown: [
    { label: 'Active SIMs (1,344)', amt: 940.8 },
    { label: 'Data egress (312 GB)', amt: 218.4 },
    { label: 'Node management (8 nodes)', amt: 96.0 },
    { label: 'Support plan', amt: 29.3 },
  ],
  revenueSeries: series(4200, 12, 0.06, 0.25),
  invoices: [
    { id: 'INV-0631', period: 'May 2026', amt: 1206.1, status: 'paid' },
    { id: 'INV-0584', period: 'Apr 2026', amt: 1142.55, status: 'paid' },
    { id: 'INV-0532', period: 'Mar 2026', amt: 1098.2, status: 'paid' },
    { id: 'INV-0489', period: 'Feb 2026', amt: 1011.75, status: 'paid' },
  ] as Invoice[],
};
