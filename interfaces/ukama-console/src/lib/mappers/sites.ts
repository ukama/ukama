/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Maps BFF composite SiteDto (+ per-site node counts) onto the existing Site
 * view-model so cards/drawers/map stay unchanged. uptime/subs/battery/signal
 * are metrics-phase data (backend gaps #7/#8) — placeholders until then.
 */
import type { ViewSiteFragment } from '@/client/graphql/views-shared.generated';
import type { Site } from '@/data';

export interface SiteNodeCounts {
  total: number;
  online: number;
}

export const toSiteStatus = (
  site: Pick<ViewSiteFragment, 'isDeactivated'>,
  counts?: SiteNodeCounts
): Site['status'] => {
  if (site.isDeactivated) return 'offline';
  if (!counts || counts.total === 0) return 'online';
  if (counts.online === 0) return 'offline';
  if (counts.online < counts.total) return 'degraded';
  return 'online';
};

export const toSite = (
  dto: ViewSiteFragment,
  counts?: SiteNodeCounts,
  subs = 0
): Site => ({
  id: dto.id,
  name: dto.name,
  area: dto.location || '—',
  status: toSiteStatus(dto, counts),
  // Network-wide subscriber count (per-site attribution is a backend gap).
  subs,
  nodes: counts?.total ?? 0,
  // TODO(metrics-phase): uptime/battery/signal/data from siteView.kpis and
  // siteView.power — backend gaps #7/#8
  uptime: 0,
  battery: 0,
  signal: null,
  data: '',
  lat: parseFloat(dto.latitude) || 0,
  lng: parseFloat(dto.longitude) || 0,
  plan: '',
});
