/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
// ESLint 9 flat config — replaces .eslintrc/.eslintignore. Mirrors the
// previous setup: eslint + @typescript-eslint recommended, prettier as a
// lint rule, plus the project's rule overrides.
import js from "@eslint/js";
import prettierRecommended from "eslint-plugin-prettier/recommended";
import globals from "globals";
import tseslint from "typescript-eslint";

export default tseslint.config(
  {
    ignores: [
      "node_modules/**",
      "dist/**",
      "build/**",
      "coverage/**",
      // parked — excluded from build/lint/run (see README)
      "planning-tool/**",
      // k6 script runs in k6's runtime (own globals like __ENV), not Node
      "load/**",
    ],
  },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  // eslint-plugin-prettier v5 "recommended": registers the plugin AND applies
  // eslint-config-prettier to disable conflicting stylistic rules.
  prettierRecommended,
  {
    languageOptions: {
      globals: { ...globals.node },
    },
    rules: {
      "no-console": "warn",
      "prettier/prettier": "error",
      "no-empty": "warn",
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/no-empty-function": "warn",
      // ban-types was removed in typescript-eslint v8; its successors:
      "@typescript-eslint/no-empty-object-type": "warn",
      "@typescript-eslint/no-unsafe-function-type": "warn",
      "@typescript-eslint/no-wrapper-object-types": "warn",
      // NOTE: no comma-dangle here — prettier owns formatting (the old
      // .eslintrc re-enabled it after eslint-config-prettier, which caused
      // circular --fix conflicts with trailingComma: "es5").
    },
  },
  {
    // Standalone worker scripts spawned by the subscriptions service; they
    // run outside the winston logger, so console IS their logging.
    files: ["threads/**"],
    rules: { "no-console": "off" },
  }
);
