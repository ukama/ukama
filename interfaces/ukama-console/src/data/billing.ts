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
