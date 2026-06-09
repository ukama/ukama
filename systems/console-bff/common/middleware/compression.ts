/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Zero-dependency gzip middleware for buffered JSON responses.
 *
 * Apollo's express4 integration (and express `res.json`) deliver complete
 * responses through `res.send`, so intercepting `send` covers all GraphQL
 * payloads — the case that matters for composite queries. Streaming/chunked
 * responses (`res.write`) intentionally pass through uncompressed.
 */
import type { NextFunction, Request, Response } from "express";
import { gzipSync } from "zlib";

/** Don't burn CPU compressing tiny payloads. */
export const COMPRESSION_THRESHOLD_BYTES = parseInt(
  process.env.COMPRESSION_THRESHOLD_BYTES ?? "1024"
);

const acceptsGzip = (req: Request): boolean => {
  const acceptEncoding = req.headers["accept-encoding"];
  return (
    typeof acceptEncoding === "string" &&
    /(^|,)\s*gzip\s*(;|,|$)/i.test(acceptEncoding)
  );
};

export const compression =
  () =>
  (req: Request, res: Response, next: NextFunction): void => {
    if (!acceptsGzip(req)) {
      next();
      return;
    }
    const originalSend = res.send.bind(res);
    res.send = (body?: unknown): Response => {
      const raw =
        typeof body === "string"
          ? Buffer.from(body)
          : Buffer.isBuffer(body)
            ? body
            : undefined;
      if (
        raw === undefined ||
        raw.length < COMPRESSION_THRESHOLD_BYTES ||
        res.getHeader("content-encoding") !== undefined
      ) {
        return originalSend(body as never);
      }
      res.setHeader("content-encoding", "gzip");
      res.setHeader("vary", "Accept-Encoding");
      const compressed = gzipSync(raw);
      res.setHeader("content-length", compressed.length);
      return originalSend(compressed);
    };
    next();
  };
