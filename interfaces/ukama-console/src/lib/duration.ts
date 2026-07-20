/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Package duration helpers.
 *
 * Since #1496 the backend stores and interprets package `duration` in
 * MINUTES (previously days). The console still presents validity in whole
 * days, so these helpers convert at the GraphQL/API boundary.
 */
export const MINUTES_PER_DAY = 1440;

/** Whole-day validity -> minutes the backend expects. */
export const daysToMinutes = (days: number): number =>
  Math.round(days) * MINUTES_PER_DAY;

/** Backend minute duration -> whole days for display. */
export const minutesToDays = (minutes: number): number =>
  Math.round(minutes / MINUTES_PER_DAY);
