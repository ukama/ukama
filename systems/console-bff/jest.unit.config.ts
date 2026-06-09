/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Unit-test config (CI gate). Runs only the dependency-free unit suite in
 * common/tests — no backend systems required — and enforces a coverage floor
 * on the security-critical helpers it exercises. The integration tests under
 * gateway/tests need live backend systems and are run separately.
 */
import type { Config } from "jest";

const config: Config = {
  preset: "ts-jest",
  testEnvironment: "jest-environment-node",
  roots: ["<rootDir>/common/tests"],
  collectCoverage: true,
  coverageProvider: "v8",
  coverageDirectory: "coverage",
  // Scope coverage to the units under test so the gate is meaningful rather
  // than diluted by the integration-heavy resolver/datasource code.
  collectCoverageFrom: [
    "common/auth/token.ts",
    "common/storage/index.ts",
    "common/configs/validateEnv.ts",
    "common/logger/requestContext.ts",
  ],
  coverageThreshold: {
    global: { lines: 60, functions: 60, statements: 60, branches: 50 },
  },
};

export default config;
