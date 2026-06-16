/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Map a site/node status to its map-pin colour (CSS variable). */
const STATUS_PIN_COLOR: Record<string, string> = {
  online: 'var(--uk-success-bright)',
  degraded: 'var(--uk-warning)',
  warning: 'var(--uk-warning)',
  offline: 'var(--uk-error)',
};

/** Pin colour for a status, falling back to the accent colour. */
export const pinColor = (status: string): string =>
  STATUS_PIN_COLOR[status] ?? 'var(--uk-ac)';
