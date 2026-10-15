/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import type {
  Fetcher,
  FetcherRequestInit,
  FetcherResponse,
} from "@apollo/utils.fetcher";

import { HTTP_TIMEOUT_MS } from "../configs";
import { logger } from "../logger";
import {
  UpstreamCall,
  getRequestId,
  recordUpstreamCall,
} from "../logger/requestContext";
import { upstreamName } from "../logger/upstreamNames";
import { REQUEST_ID_HEADER } from "../middleware/requestId";

/** Host and path for the log line; never throws, so fetch reports bad URLs as before. */
const target = (
  url: string
): { host: string; service?: string; path: string } => {
  try {
    const { host, pathname, search } = new URL(url);
    const service = upstreamName(host);
    return { host, ...(service && { service }), path: pathname + search };
  } catch {
    return { host: "", path: url };
  }
};

const describeFailure = (err: unknown, timedOut: boolean, limitMs: number) => {
  if (timedOut) return `timeout after ${limitMs}ms`;
  const cause = (err as { cause?: { code?: string } })?.cause?.code;
  return cause ?? (err instanceof Error ? err.message : String(err));
};

/**
 * One line per upstream call: debug when it succeeds, info for 4xx (the
 * resolver surfaces those), warn for 5xx and calls that got no response.
 */
const logUpstreamCall = (call: UpstreamCall): void => {
  recordUpstreamCall(call);
  const level =
    call.failure || call.status >= 500
      ? "warn"
      : call.status >= 400
        ? "info"
        : "debug";
  logger.log(
    level,
    `upstream ${call.method} ${call.service ? `${call.service} ` : call.host}${call.path} ${call.failure ?? call.status} ${call.durationMs}ms`,
    { upstream: call }
  );
};

/**
 * Sends one upstream call with a hard time limit, forwards the request's
 * correlation id, and logs the outcome.
 */
const timedFetch = async (
  url: string,
  init: FetcherRequestInit | undefined,
  timeoutMs: number
): Promise<FetcherResponse> => {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeoutMs);
  const requestId = getRequestId();
  const callerSetId = Object.keys(init?.headers ?? {}).some(
    h => h.toLowerCase() === REQUEST_ID_HEADER
  );
  const headers =
    requestId && !callerSetId
      ? { ...(init?.headers ?? {}), [REQUEST_ID_HEADER]: requestId }
      : init?.headers;
  const call = {
    method: (init?.method ?? "GET").toUpperCase(),
    ...target(url),
  };
  const started = Date.now();
  try {
    const response = await fetch(url, {
      ...(init ?? {}),
      headers,
      signal: controller.signal,
    } as RequestInit);
    logUpstreamCall({
      ...call,
      status: response.status,
      durationMs: Date.now() - started,
    });
    // Node's fetch Response satisfies FetcherResponse at runtime; the cast is
    // needed because tsconfig's "dom" lib (without "dom.iterable") types
    // Headers without its iterator methods, which trips ts-node's checker.
    return response as unknown as FetcherResponse;
  } catch (err) {
    logUpstreamCall({
      ...call,
      status: 0,
      durationMs: Date.now() - started,
      failure: describeFailure(err, controller.signal.aborted, timeoutMs),
    });
    throw err;
  } finally {
    clearTimeout(timer);
  }
};

/**
 * A fetch wrapper that aborts any upstream HTTP call exceeding
 * HTTP_TIMEOUT_MS, so a slow backend can never hang the BFF.
 */
const fetchWithTimeout: Fetcher = (url, init) =>
  timedFetch(url, init, HTTP_TIMEOUT_MS);

/** Same as fetchWithTimeout, with a caller-chosen limit. */
export const createTimedFetch =
  (timeoutMs: number): Fetcher =>
  (url, init) =>
    timedFetch(url, init, timeoutMs);

/**
 * For the few upstream calls that legitimately outlive HTTP_TIMEOUT_MS.
 */
export class TimedRESTDataSource extends RESTDataSource {
  constructor(timeoutMs: number) {
    super({ fetch: createTimedFetch(timeoutMs), logger });
  }
}

/**
 * Base class for all subgraph datasources. Extends Apollo's
 * RESTDataSource with a hard request timeout. All datasources
 * must extend this class instead of RESTDataSource directly.
 */
export class BaseRESTDataSource extends RESTDataSource {
  constructor() {
    super({ fetch: fetchWithTimeout, logger });
  }
}
