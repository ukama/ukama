/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

type LogLevel = 'info' | 'warn' | 'error';

interface LogEvent {
  level: LogLevel;
  message: string;
  context?: Record<string, unknown>;
  error?: unknown;
}

const log = (event: LogEvent): void => {
  if (process.env.NODE_ENV === 'development') {
    const { level, message, context, error } = event;
    const fn = console[level] ?? console.log;
    fn(`[${level.toUpperCase()}] ${message}`, context ?? '', error ?? '');
  }
  // Future: send to Sentry / LogRocket here
};

export const logger = {
  info: (message: string, context?: Record<string, unknown>) =>
    log({ level: 'info', message, context }),
  warn: (message: string, context?: Record<string, unknown>) =>
    log({ level: 'warn', message, context }),
  error: (message: string, error?: unknown, context?: Record<string, unknown>) =>
    log({ level: 'error', message, error, context }),
};
