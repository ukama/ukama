/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Time-range selector shared by the metric charts (Day/Week/Month). */
export type Range = 'Day' | 'Week' | 'Month';

export const RANGES: Range[] = ['Day', 'Week', 'Month'];

/** Window length per range, in seconds (drives the metricsRange from/to). */
export const RANGE_SECONDS: Record<Range, number> = {
  Day: 86_400,
  Week: 604_800,
  Month: 2_592_000,
};
