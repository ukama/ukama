/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Shared types/colors for the business sites map. */

export interface SiteMapSite {
  id: string;
  name: string;
  status: 'online' | 'warning' | 'degraded' | 'offline';
  lat: number;
  lng: number;
}

export const BIZ_DOT: Record<SiteMapSite['status'], string> = {
  online: 'var(--uk-success)',
  warning: 'var(--uk-warning)',
  degraded: 'var(--uk-warning)',
  offline: 'var(--uk-error)',
};
