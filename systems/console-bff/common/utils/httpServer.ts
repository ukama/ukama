/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import type { Server } from "http";

import { HTTP_KEEP_ALIVE_TIMEOUT_MS } from "../configs";

/**
 * Holds idle keep-alive sockets open longer than the proxy in front of the
 * server, so the proxy always closes first and never reuses a dead socket.
 */
export const applyKeepAliveTimeouts = (
  server: Server,
  keepAliveTimeoutMs: number = HTTP_KEEP_ALIVE_TIMEOUT_MS
): Server => {
  server.keepAliveTimeout = keepAliveTimeoutMs;
  // Node's default headersTimeout (60s) is below the keep-alive, so raise it too.
  server.headersTimeout = keepAliveTimeoutMs + 1000;
  return server;
};
