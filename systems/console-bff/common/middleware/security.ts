/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Security middleware for the gateway: helmet for response headers.
 */
import helmet from "helmet";

/**
 * Security response headers. CSP is disabled because this is a JSON API
 * gateway (no first-party HTML to protect); the frontend enforces its own
 * CSP. All other helmet protections (nosniff, frameguard, HSTS, etc.) apply.
 */
export const securityHeaders = () => helmet({ contentSecurityPolicy: false });
