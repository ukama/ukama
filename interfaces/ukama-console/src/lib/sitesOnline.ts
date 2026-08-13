/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * "Sites online" tile value (Ops + Business home). The count comes from the
 * analytics SITES_ONLINE KPI and the total from the registry `sitesView` —
 * two responses that are fetched and expired separately, so the pair is only
 * shown when both sides are known and `online <= total` holds.
 */

export interface SitesOnlineTile {
  /** Online count — set only when the pair is coherent. */
  online?: number;
  /** Registry site total — set only when the pair is coherent. */
  total?: number;
  /** Sites not fully online — set only when the pair is coherent. */
  offline?: number;
  /** Display value: "1/3", or "—". */
  text: string;
}

const UNKNOWN: SitesOnlineTile = { text: '—' };

/** Pass `undefined` for a side that is unknown; `0` means genuinely zero. */
export const sitesOnlineTile = (
  onlineKpi: number | undefined,
  total: number | undefined,
): SitesOnlineTile => {
  const online = onlineKpi != null ? Math.round(onlineKpi) : undefined;
  if (online == null || total == null || online > total) return UNKNOWN;
  return { online, total, offline: total - online, text: `${online}/${total}` };
};
