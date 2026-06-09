/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Per-request service base-URL resolution for composite resolvers.
 *
 * NOTE: the consolidated AppContext does not carry `ctx.baseURL` (each
 * domain's legacy Context interface declares it, but the runtime context
 * never sets it). Composites therefore resolve base URLs explicitly — same
 * pattern as getOrgTree — memoized per request so multiple sections hitting
 * the same upstream system cost one lookup.
 */
import { getBaseURL } from "../common/utils";

export class ServiceUrlResolver {
  private readonly cache = new Map<string, Promise<string>>();

  constructor(private readonly orgName: string) {}

  /** Resolve (and memoize) the API-gateway base URL for a SUB_GRAPHS service
   *  name (e.g. "network", "sim", "health"). Throws on lookup failure —
   *  callers run inside runSection, which maps it to a SectionError. */
  url(service: string): Promise<string> {
    let cached = this.cache.get(service);
    if (!cached) {
      cached = getBaseURL(service, this.orgName).then(res => {
        if (res.status !== 200 || !res.message) {
          // getBaseURL's message now carries the diagnosis (empty orgName,
          // missing system mapping, or the init fetch error + target URL).
          throw new Error(
            `Base URL lookup failed for service '${service}': ${res.message}`
          );
        }
        return res.message;
      });
      this.cache.set(service, cached);
    }
    return cached;
  }
}
