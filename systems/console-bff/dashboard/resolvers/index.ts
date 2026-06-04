/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { NetworkOverviewResolver } from "./networkOverview";
import { NodeViewResolver, NodesViewResolver } from "./nodeViews";
import { SimPoolViewResolver } from "./simPoolView";
import { SiteViewResolver, SitesViewResolver } from "./siteViews";
import { SubscribersViewResolver } from "./subscribersView";

/** View-domain composite resolvers (plan §3) — network lens, Phase 2. */
const dashboardResolvers = [
  NetworkOverviewResolver,
  NodesViewResolver,
  NodeViewResolver,
  SitesViewResolver,
  SiteViewResolver,
  SubscribersViewResolver,
  SimPoolViewResolver,
];

export default dashboardResolvers;
