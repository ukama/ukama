/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Apollo plugin that writes one log line per GraphQL operation: name,
 * duration, caller, error codes and the upstream calls it made. GraphQL
 * failures come back as HTTP 200, so the access log alone cannot show them.
 */
import type { ApolloServerPlugin, BaseContext } from "@apollo/server";
import type { GraphQLError } from "graphql";
import type { Logger } from "winston";

import { logger as defaultLogger } from ".";
import { UpstreamCall, getUpstreamCalls } from "./requestContext";

// Enough to diagnose a failure without letting one bad query flood the line.
const MAX_LOGGED_ERRORS = 5;

interface CallerContext extends BaseContext {
  headers?: { userId?: string; orgName?: string };
}

export interface LoggedError {
  code: string;
  message: string;
  path?: string;
  upstreamStatus?: number;
}

export interface UpstreamSummary {
  method: string;
  host: string;
  result: string;
  count: number;
  maxMs: number;
}

/**
 * Apollo's own code when it set one, else the upstream HTTP status that
 * RESTDataSource attached, else a generic internal error.
 */
export const errorCode = (err: GraphQLError): string => {
  const code = err.extensions?.code;
  if (typeof code === "string") return code;
  const status = (err.extensions?.response as { status?: number } | undefined)
    ?.status;
  return status ? `UPSTREAM_${status}` : "INTERNAL_SERVER_ERROR";
};

const toLoggedError = (err: GraphQLError): LoggedError => {
  const upstreamStatus = (
    err.extensions?.response as { status?: number } | undefined
  )?.status;
  return {
    code: errorCode(err),
    message: err.message,
    ...(err.path && { path: err.path.join(".") }),
    ...(upstreamStatus && { upstreamStatus }),
  };
};

/** Collapses repeated calls (e.g. seven init lookups) into one entry each. */
export const summariseUpstream = (
  calls: readonly UpstreamCall[]
): UpstreamSummary[] => {
  const byKey = new Map<string, UpstreamSummary>();
  for (const call of calls) {
    const result = call.failure ?? String(call.status);
    const key = `${call.method} ${call.host} ${result}`;
    const entry = byKey.get(key);
    if (entry) {
      entry.count += 1;
      entry.maxMs = Math.max(entry.maxMs, call.durationMs);
    } else {
      byKey.set(key, {
        method: call.method,
        host: call.host,
        result,
        count: 1,
        maxMs: call.durationMs,
      });
    }
  }
  return [...byKey.values()];
};

const isFailedUpstream = (u: UpstreamSummary) => !/^[23]\d\d$/.test(u.result);

export const operationLoggerPlugin = (
  log: Logger = defaultLogger
): ApolloServerPlugin<CallerContext> => ({
  async requestDidStart() {
    const started = Date.now();
    return {
      async willSendResponse(ctx) {
        // Logging must never turn a good response into a failed one.
        try {
          const op =
            ctx.operationName ?? ctx.request.operationName ?? "anonymous";
          const durationMs = Date.now() - started;
          const status = ctx.response.http.status ?? 200;
          const errors = (ctx.errors ?? [])
            .slice(0, MAX_LOGGED_ERRORS)
            .map(toLoggedError);
          const upstream = summariseUpstream(getUpstreamCalls());
          const failed = upstream.filter(isFailedUpstream);

          let message = `graphql ${op} ${durationMs}ms`;
          if (errors.length) {
            message += ` errors=${[...new Set(errors.map(e => e.code))].join(",")}`;
          }
          if (failed.length) {
            message += ` upstream=${failed.map(u => `${u.host} ${u.result}`).join(",")}`;
          }

          const level =
            status >= 500 ? "error" : errors.length ? "warn" : "info";
          log.log(level, message, {
            op,
            opType: ctx.operation?.operation,
            durationMs,
            status,
            userId: ctx.contextValue?.headers?.userId || undefined,
            org: ctx.contextValue?.headers?.orgName || undefined,
            ...(errors.length && {
              errors,
              errorCount: ctx.errors?.length,
            }),
            ...(upstream.length && { upstream }),
          });
        } catch (err) {
          log.debug(`operation logger failed: ${err}`);
        }
      },
    };
  },
});
