/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Business-lens view-model types. The seed datasets that used to live here
 * were removed once the Business screens were wired to the analytics service;
 * only the shared `BizSite` shape (used by the Home sites map) remains.
 */

export interface BizSite {
  id: string;
  name: string;
  status: 'online' | 'warning' | 'offline';
  revenue: number;
  revToday: number;
  customers: number;
  custToday: number;
  data: string;
  uptime: number;
  top: string;
  issue: string | null;
  lat: number;
  lng: number;
}
