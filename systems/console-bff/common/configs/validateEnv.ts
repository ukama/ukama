/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Startup environment validation. Fails fast with a clear message when a
 * required variable is missing, so misconfiguration is caught at boot rather
 * than at the first request. In production it additionally rejects localhost
 * defaults, which almost always indicate a broken deployment.
 */
import { logger } from "../logger";
import {
  AUTH_URL,
  INIT_API_GW,
  INVENTORY_API_GW,
  IS_PRODUCTION,
  NUCLEUS_API_GW,
} from "./index";

interface EnvVar {
  name: string;
  value: string | undefined;
  /** A URL whose value must not point at localhost in production. */
  url?: boolean;
}

/**
 * Validates required environment. Throws on the first hard failure.
 * `JWT_SECRET` is validated separately where the token secret is derived.
 */
export const validateEnv = (): void => {
  const required: EnvVar[] = [
    { name: "AUTH_URL", value: AUTH_URL, url: true },
    { name: "NUCLEUS_API_GW", value: NUCLEUS_API_GW, url: true },
    { name: "INIT_API_GW", value: INIT_API_GW, url: true },
    { name: "INVENTORY_API_GW", value: INVENTORY_API_GW, url: true },
  ];

  const errors: string[] = [];

  for (const { name, value, url } of required) {
    if (!value || value.trim() === "") {
      errors.push(`${name} is required but not set`);
      continue;
    }
    if (url) {
      try {
        new URL(value);
      } catch {
        errors.push(`${name} is not a valid URL: "${value}"`);
        continue;
      }
      if (IS_PRODUCTION && /localhost|127\.0\.0\.1/.test(value)) {
        errors.push(`${name} points at localhost in production ("${value}")`);
      }
    }
  }

  if (errors.length > 0) {
    const message = `Invalid environment configuration:\n - ${errors.join(
      "\n - "
    )}`;
    logger.error(message);
    throw new Error(message);
  }

  logger.info("Environment validation passed");
};
