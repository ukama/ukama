/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Security middleware for the gateway: helmet for response headers and
 * express-rate-limit for per-client throttling.
 */
import { rateLimit as expressRateLimit } from "express-rate-limit";
import helmet from "helmet";

const RATE_LIMIT_WINDOW_MS = 60_000;
const RATE_LIMIT_MAX = 300;
const RATE_LIMIT_SKIP_PATHS = new Set(["/healthz", "/readyz", "/ping"]);

/**
 * Security response headers. CSP is disabled because this is a JSON API
 * gateway (no first-party HTML to protect); the frontend enforces its own
 * CSP. All other helmet protections (nosniff, frameguard, HSTS, etc.) apply.
 */
export const securityHeaders = () => helmet({ contentSecurityPolicy: false });

/** Fixed-window per-IP rate limiter; health/ping paths are exempt. */
export const rateLimit = () =>
  expressRateLimit({
    windowMs: RATE_LIMIT_WINDOW_MS,
    limit: RATE_LIMIT_MAX,
    standardHeaders: "draft-7",
    legacyHeaders: false,
    skip: req => RATE_LIMIT_SKIP_PATHS.has(req.path),
  });
