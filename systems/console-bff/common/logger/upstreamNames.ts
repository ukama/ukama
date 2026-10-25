/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Maps upstream hosts to system names for log lines. In-cluster systems are
 * reached by ClusterIP (e.g. 172.20.188.168:8080), which says nothing on its
 * own; getBaseURL records each one as init resolves it.
 */
import {
  AUTH_URL,
  INIT_API_GW,
  INVENTORY_API_GW,
  NUCLEUS_API_GW,
} from "../configs";

const names = new Map<string, string>();

/** Remembers which system `baseURL` belongs to. Ignores unparseable URLs. */
export const nameUpstream = (baseURL: string, system: string): void => {
  if (!baseURL || !system) return;
  try {
    names.set(new URL(baseURL).host, system);
  } catch {
    // A bad URL fails the upstream call itself; the log label can stay unknown.
  }
};

/** System name for a host, if the BFF has resolved or configured it. */
export const upstreamName = (host: string): string | undefined =>
  names.get(host);

nameUpstream(INIT_API_GW, "init");
nameUpstream(NUCLEUS_API_GW, "nucleus");
nameUpstream(INVENTORY_API_GW, "inventory");
nameUpstream(AUTH_URL, "auth");
