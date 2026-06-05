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
 * Logger with two output formats:
 *
 * - "json"  — structured JSON, one object per line (timestamp, level,
 *             message, metadata) — ready for log aggregators. Default in
 *             production.
 * - "pretty" — human-readable, level-colored lines for local development.
 *             Muted ("matte") ANSI-256 colors: red for error, yellow for
 *             warn, blue for info; everything else stays default/white.
 *             Default outside production when stdout is a TTY.
 *
 * Override with LOG_FORMAT=json|pretty. LOG_LEVEL controls verbosity
 * (default "info").
 */

/** Muted ANSI-256 foreground colors per level (whole line is tinted). */
const LEVEL_COLOR: Record<string, string> = {
  error: "\x1b[38;5;167m", // matte red
  warn: "\x1b[38;5;179m", // matte yellow
  info: "\x1b[38;5;110m", // matte blue
};
const RESET = "\x1b[0m";
const DIM = "\x1b[2m";

/** Stamps every log line with the current request's correlation id, if set. */
const withRequestId = winston.format(info => {
  const requestId = getRequestId();
  if (requestId) info.requestId = requestId;
  return info;
});

const prettyPrint = winston.format.printf(info => {
  const { timestamp, level, message, stack, ...meta } = info as {
    timestamp?: string;
    level: string;
    message: unknown;
    stack?: string;
    [key: string]: unknown;
  };
  delete meta[Symbol.for("level") as unknown as string];

  const color = LEVEL_COLOR[level] ?? "";
  const time = timestamp ? `${DIM}${timestamp}${RESET} ` : "";
  const label = level.toUpperCase().padEnd(5);
  const requestId = meta.requestId ? `${DIM}[${meta.requestId}]${RESET} ` : "";
  delete meta.requestId;
  delete meta.service;

  const extra = Object.keys(meta).length
    ? ` ${DIM}${JSON.stringify(meta)}${RESET}`
    : "";
  const trace = stack ? `\n${stack}` : "";

  const body = `${label} ${message}${trace}`;
  // Tint level + message per level; timestamp/requestId/meta stay dim.
  const colored = color ? `${color}${body}${RESET}` : body;
  return `${time}${requestId}${colored}${extra}`;
});

const useJson =
  process.env.LOG_FORMAT === "json" ||
  (process.env.LOG_FORMAT !== "pretty" &&
    (process.env.NODE_ENV === "production" || !process.stdout.isTTY));

const format = useJson
  ? winston.format.combine(
      withRequestId(),
      winston.format.timestamp(),
      winston.format.errors({ stack: true }),
      winston.format.json()
    )
  : winston.format.combine(
      withRequestId(),
      winston.format.timestamp({ format: "HH:mm:ss" }),
      winston.format.errors({ stack: true }),
      prettyPrint
    );

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL ?? "info",
  defaultMeta: process.env.SERVICE_NAME
    ? { service: process.env.SERVICE_NAME }
    : undefined,
  transports: [new winston.transports.Console()],
  format,
});

export { logger };
