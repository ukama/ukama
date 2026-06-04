/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Persisted-operations allowlist (docs/bff-screen-api-plan.md §3.3).
 *
 * The console is the BFF's only client, so production accepts only the
 * operation documents the console actually ships. The console build emits a
 * manifest (graphql-codegen persisted-documents format:
 * `{ "<sha256-of-document>": "<document>" }`); this middleware hashes each
 * incoming query and rejects unknown documents.
 *
 * Disabled unless PERSISTED_OPS_ENFORCED=true (dev and codegen stay open).
 * Introspection is governed separately by INTROSPECTION_ENABLED.
 */
import { createHash } from "crypto";
import type { NextFunction, Request, Response } from "express";
import { existsSync, readFileSync } from "fs";

import { logger } from "../logger";

export const PERSISTED_OPS_ENFORCED =
  process.env.PERSISTED_OPS_ENFORCED === "true";
export const PERSISTED_OPS_MANIFEST =
  process.env.PERSISTED_OPS_MANIFEST ?? "persisted-operations.json";

export const sha256 = (document: string): string =>
  createHash("sha256").update(document).digest("hex");

export const loadManifest = (path: string): Set<string> => {
  if (!existsSync(path)) {
    throw new Error(`Persisted operations manifest not found: ${path}`);
  }
  const manifest = JSON.parse(readFileSync(path, "utf8")) as Record<
    string,
    string
  >;
  return new Set(Object.keys(manifest));
};

const isIntrospectionOperation = (query: unknown): boolean =>
  typeof query === "string" && query.includes("__schema");

/**
 * Express middleware guarding /graphql. Mount after body parsing, before
 * the Apollo middleware. No-op factory result when enforcement is off.
 */
export const persistedOperations = (
  options: {
    enforced?: boolean;
    manifestPath?: string;
    allowIntrospection?: boolean;
  } = {}
) => {
  const enforced = options.enforced ?? PERSISTED_OPS_ENFORCED;
  const allowIntrospection = options.allowIntrospection ?? false;
  if (!enforced) {
    return (_req: Request, _res: Response, next: NextFunction) => next();
  }

  // Fail fast at startup if enforcement is on but the manifest is missing.
  const allowedHashes = loadManifest(
    options.manifestPath ?? PERSISTED_OPS_MANIFEST
  );
  logger.info(
    `Persisted operations enforced: ${allowedHashes.size} known operations`
  );

  return (req: Request, res: Response, next: NextFunction): void => {
    const ops = Array.isArray(req.body) ? req.body : [req.body];
    for (const op of ops) {
      const query = (op as { query?: unknown })?.query;
      if (allowIntrospection && isIntrospectionOperation(query)) {
        continue;
      }
      if (typeof query !== "string" || !allowedHashes.has(sha256(query))) {
        logger.warn("Rejected non-persisted GraphQL operation", {
          operationName: (op as { operationName?: string })?.operationName,
        });
        res.status(403).json({
          errors: [
            {
              message: "Operation not allowed",
              extensions: { code: "PERSISTED_OPERATION_NOT_FOUND" },
            },
          ],
        });
        return;
      }
    }
    next();
  };
};
