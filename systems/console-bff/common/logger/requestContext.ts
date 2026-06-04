/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Per-request context backed by AsyncLocalStorage so the correlation id flows
 * implicitly through async calls and every log line can pick it up — no need
 * to thread a requestId argument through resolvers and datasources.
 *
 *  - Express servers (the gateway) use `runWithRequestId` to wrap the whole
 *    request chain.
 *  - Apollo standalone servers (the subgraphs) can't add middleware, so they
 *    call `setRequestId` from the per-request context function; it sets the
 *    store for the remainder of that request's async execution.
 */
import { AsyncLocalStorage } from "async_hooks";

interface RequestContext {
  requestId: string;
}

const storage = new AsyncLocalStorage<RequestContext>();

/** Returns the current request's correlation id, if any. */
export const getRequestId = (): string | undefined =>
  storage.getStore()?.requestId;

/** Runs `fn` within a request context (preferred — fully scoped). */
export const runWithRequestId = <T>(requestId: string, fn: () => T): T =>
  storage.run({ requestId }, fn);

/**
 * Binds a request id to the current async execution without a wrapping
 * callback. Used where middleware isn't available (Apollo context functions).
 */
export const setRequestId = (requestId: string): void => {
  if (!requestId) return;
  storage.enterWith({ requestId });
};
