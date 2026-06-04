/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { logger } from "../common/logger";
import { SectionError, SectionErrorCode } from "./types";

/** Per-section upstream timeout — tighter than the global HTTP_TIMEOUT_MS
 *  so one slow upstream degrades a section, not the whole screen. */
export const SECTION_TIMEOUT_MS = parseInt(
  process.env.SECTION_TIMEOUT_MS ?? "5000"
);

/**
 * Request-scoped collector for composite section failures. One instance is
 * created per composite root resolver and threaded to its section resolvers;
 * the composite's `errors` field resolves `list()` last.
 */
export class SectionErrorCollector {
  private readonly errors: SectionError[] = [];

  add(section: string, code: SectionErrorCode, message: string): void {
    this.errors.push({ section, code, message });
  }

  list(): SectionError[] {
    return this.errors;
  }
}

/**
 * Result of one composite section: the value (null on failure) plus the
 * typed error (null on success). Sections embed `error` in their own GraphQL
 * type rather than a composite-level errors list because GraphQL executes
 * sibling field resolvers concurrently — a top-level `errors` field would
 * resolve before the sections it reports on.
 */
export interface SectionResult<T> {
  value: T | null;
  error: SectionError | null;
}

/**
 * Like `withSection`, but returns the error alongside the value so section
 * field resolvers can embed it in their section type (the composite-query
 * pattern). Same timeout + logging behavior.
 */
export async function runSection<T>(
  section: string,
  fn: () => Promise<T>,
  timeoutMs: number = SECTION_TIMEOUT_MS
): Promise<SectionResult<T>> {
  const collector = new SectionErrorCollector();
  const value = await withSection(collector, section, fn, timeoutMs);
  return { value, error: collector.list()[0] ?? null };
}

/** `SectionResult` for a backend gap (§3.5): null value + NOT_IMPLEMENTED. */
export const notImplementedSection = <T>(
  section: string
): SectionResult<T> => ({
  value: null,
  error: {
    section,
    code: SectionErrorCode.NOT_IMPLEMENTED,
    message: `${section} is not available yet`,
  },
});

/**
 * Marks a section whose backend endpoint/property doesn't exist yet
 * (docs/bff-screen-api-plan.md §3.5). Call from the section resolver next to
 * its `TODO(backend-gap)` comment and return null. Never fabricate data for
 * a gap — the console renders a placeholder off this code.
 */
export const sectionNotImplemented = (
  collector: SectionErrorCollector,
  section: string
): null => {
  collector.add(
    section,
    SectionErrorCode.NOT_IMPLEMENTED,
    `${section} is not available yet`
  );
  return null;
};

const toSectionErrorCode = (error: unknown): SectionErrorCode => {
  const err = error as {
    message?: string;
    code?: number | string;
    extensions?: { response?: { status?: number } };
  };
  const status =
    err?.extensions?.response?.status ??
    (typeof err?.code === "number" ? err.code : undefined);
  if (err?.message === "SECTION_TIMEOUT") {
    return SectionErrorCode.UPSTREAM_TIMEOUT;
  }
  if (status === 404) return SectionErrorCode.NOT_FOUND;
  if (status === 401 || status === 403) return SectionErrorCode.FORBIDDEN;
  if (typeof status === "number") return SectionErrorCode.UPSTREAM_ERROR;
  return SectionErrorCode.INTERNAL;
};

/**
 * Errors-as-data wrapper for composite section resolvers
 * (docs/bff-screen-api-plan.md §4.5). Every section resolver MUST run its
 * upstream work through this wrapper — never hand-roll try/catch.
 *
 * On success: returns the value and logs section latency (requestId is
 * stamped by the logger's AsyncLocalStorage context).
 * On failure/timeout: logs, records a typed SectionError on the collector,
 * and returns null so the rest of the composite still resolves.
 */
export async function withSection<T>(
  collector: SectionErrorCollector,
  section: string,
  fn: () => Promise<T>,
  timeoutMs: number = SECTION_TIMEOUT_MS
): Promise<T | null> {
  const startedAt = Date.now();
  let timer: NodeJS.Timeout | undefined;
  try {
    const result = await Promise.race([
      fn(),
      new Promise<never>((_, reject) => {
        timer = setTimeout(
          () => reject(new Error("SECTION_TIMEOUT")),
          timeoutMs
        );
      }),
    ]);
    logger.info("composite section resolved", {
      section,
      durationMs: Date.now() - startedAt,
    });
    return result;
  } catch (error) {
    const code = toSectionErrorCode(error);
    logger.error("composite section failed", {
      section,
      code,
      durationMs: Date.now() - startedAt,
      error: error instanceof Error ? error.message : String(error),
    });
    collector.add(section, code, `Failed to load ${section}`);
    return null;
  } finally {
    if (timer) clearTimeout(timer);
  }
}
