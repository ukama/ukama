/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Single source of truth for static/seed data (BUILD-PLAN §13.5).
 * Entity datasets are added here as screens are built (plans, members,
 * sims, billing, business…).
 */
export * from './networks';
export * from './sites';
export * from './nodes';
export * from './subscribers';
export * from './alerts';
export * from './reference/statuses';
