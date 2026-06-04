/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Liveness/readiness probe path used by the console helm chart (legacy parity). */
export function GET() {
  return new Response('pong', { status: 200 });
}
