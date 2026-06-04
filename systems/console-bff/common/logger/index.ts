/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import winston from "winston";

import { getRequestId } from "./requestContext";

/**
 * Structured JSON logger. One JSON object per line with a timestamp, level,
 * message and any attached metadata (e.g. requestId) — ready for ingestion by
 * a log aggregator. Set LOG_LEVEL to control verbosity (default "info").
 */

/** Stamps every log line with the current request's correlation id, if set. */
const withRequestId = winston.format(info => {
  const requestId = getRequestId();
  if (requestId) info.requestId = requestId;
  return info;
});

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL ?? "info",
  defaultMeta: process.env.SERVICE_NAME
    ? { service: process.env.SERVICE_NAME }
    : undefined,
  transports: [new winston.transports.Console()],
  format: winston.format.combine(
    withRequestId(),
    winston.format.timestamp(),
    winston.format.errors({ stack: true }),
    winston.format.json()
  ),
});

export { logger };
