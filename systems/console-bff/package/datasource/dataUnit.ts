/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Data-unit mapping between the console and the data-plan backend.
 *
 * The backend persists the data unit as an enum:
 *   "unknown" | "Bytes" | "KiloBytes" | "MegaBytes" | "GigaBytes"
 * The console works with short labels ("Bytes" | "KB" | "MB" | "GB"). This
 * module is the single place that converts between the two so the wire
 * conventions stay consistent in both directions (write on addPackage, read on
 * get/getPackages mappers).
 */

/** Backend enum values for a package's data unit. */
export type BackendDataUnit =
  | "unknown"
  | "Bytes"
  | "KiloBytes"
  | "MegaBytes"
  | "GigaBytes";

/** Short label the console sends/expects. */
export type ConsoleDataUnit = "Bytes" | "KB" | "MB" | "GB";

const CONSOLE_TO_BACKEND: Record<string, BackendDataUnit> = {
  bytes: "Bytes",
  kb: "KiloBytes",
  kilobytes: "KiloBytes",
  mb: "MegaBytes",
  megabytes: "MegaBytes",
  gb: "GigaBytes",
  gigabytes: "GigaBytes",
};

const BACKEND_TO_CONSOLE: Record<string, ConsoleDataUnit> = {
  bytes: "Bytes",
  kilobytes: "KB",
  megabytes: "MB",
  gigabytes: "GB",
};

/** Console label → backend enum (defaults to "unknown" for anything else). */
export const toBackendDataUnit = (unit?: string): BackendDataUnit =>
  CONSOLE_TO_BACKEND[(unit ?? "").trim().toLowerCase()] ?? "unknown";

/**
 * Backend enum → console label. Accepts already-short labels too (idempotent),
 * defaults to "MB" for unknown/unrecognised values so the UI always has a unit.
 */
export const toConsoleDataUnit = (unit?: string): ConsoleDataUnit => {
  const key = (unit ?? "").trim().toLowerCase();
  return (
    BACKEND_TO_CONSOLE[key] ??
    (["bytes", "kb", "mb", "gb"].includes(key)
      ? (key.toUpperCase() as ConsoleDataUnit)
      : "MB")
  );
};
