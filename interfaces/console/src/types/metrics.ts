/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Node, SiteDto as Site } from '@/client/graphql/generated';
import {
  Graphs_Type,
  LatestMetricSubRes,
} from '@/client/graphql/generated/subscriptions';

interface MetricsSubData {
  getMetricByTabSub: LatestMetricSubRes;
  getMetricStatSub: LatestMetricSubRes;
  getSiteMetricByTabSub: LatestMetricSubRes;
  getSiteMetricStatSub: LatestMetricSubRes;
  getMetricBySiteSub: LatestMetricSubRes;
}

export interface TMetricResDto {
  data: MetricsSubData;
}

export type KPIType =
  | 'node'
  | 'solar'
  | 'controller'
  | 'battery'
  | 'backhaul'
  | 'switch';

export interface ActiveView {
  graphType: Graphs_Type;
  kpi: KPIType;
}

export type TStatusBarObj = Node | Site;
