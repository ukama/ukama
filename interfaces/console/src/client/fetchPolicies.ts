/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Standardised Apollo fetchPolicy values for the Ukama Console.
 * Always pick from this list — never hardcode strings in page files.
 */
import type { WatchQueryFetchPolicy } from '@apollo/client';

/** Data that changes frequently (nodes, sites, metrics) — show cached then refresh */
export const LIVE_DATA: WatchQueryFetchPolicy = 'cache-and-network';

/** Reference data that rarely changes (currencies, packages, SIM types) */
export const STATIC_DATA: WatchQueryFetchPolicy = 'cache-first';

/** User-triggered actions where stale data is unacceptable (billing, member list) */
export const FRESH_DATA: WatchQueryFetchPolicy = 'network-only';
