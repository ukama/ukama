/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Dashboard module: view-domain composite queries for ukama-console
 * (docs/bff-screen-api-plan.md §3). Phase 1 ships the shared error layer;
 * composite resolvers land in phase 2+.
 */
export {
  SECTION_TIMEOUT_MS,
  SectionErrorCollector,
  sectionNotImplemented,
  withSection,
} from "./section";
export { SectionError, SectionErrorCode, UserError } from "./types";
