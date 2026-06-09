/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Fabricates a plausible time series (prototype data.jsx helper). */
export const series = (
  base: number,
  n = 14,
  jitter = 0.12,
  trend = 0,
): number[] =>
  Array.from({ length: n }, (_, i) =>
    Math.max(
      0,
      Math.round(
        (base +
          base * trend * (i / n) +
          base * jitter * (Math.sin(i * 1.7) + Math.cos(i * 0.9))) *
          10,
      ) / 10,
    ),
  );
