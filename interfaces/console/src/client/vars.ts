/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { makeVar } from '@apollo/client';
import { Graphs_Type, MetricsRes } from './graphql/generated/subscriptions';

export const activeGraphTypeVar = makeVar<Graphs_Type>(Graphs_Type.Home);

export const activeNodeTabVar = makeVar<number>(0);

// Reactive var for real-time node metrics (replaces PubSub → useState chain)
export const nodeMetricsVar = makeVar<MetricsRes>({ metrics: [] });

// Reactive var for site active subscriber count (replaces PubSub → useState chain)
export const siteActiveSubscribersVar = makeVar<number>(0);
