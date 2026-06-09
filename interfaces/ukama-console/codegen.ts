/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * GraphQL codegen — near-operation-file (plan §15.5): hooks/types are
 * generated NEXT TO each .graphql document (<name>.generated.ts), with
 * shared schema types in src/client/graphql/types.ts.
 *
 * Gateway schema only for now (polling-only v1). metrics.graphql is kept
 * as reference and excluded — the metrics endpoint gets its own client +
 * codegen target in the metrics phase.
 *
 * Run with the gateway up: `pnpm codegen`.
 */
import type { CodegenConfig } from '@graphql-codegen/cli';

const GW = process.env.NEXT_PUBLIC_API_GW ?? 'http://localhost:8080';

/**
 * Schema endpoint for introspection. The preflight header satisfies Apollo
 * Server's CSRF prevention (same header the Apollo Sandbox sends); no auth is
 * needed because the gateway lets introspection past the auth gate when
 * ENABLE_INTROSPECTION is set.
 */
const schema = process.env.SCHEMA_PATH
  ? // Offline mode: SDL file exported from the BFF (SCHEMA_PATH=schema.graphql
    // pnpm codegen) — used by CI's contract check and when the gateway isn't
    // running locally.
    process.env.SCHEMA_PATH
  : {
      [`${GW}/graphql`]: {
        headers: { 'apollo-require-preflight': 'true' },
      },
    };

const hooksConfig = {
  withHooks: true,
  // ergonomic generated names: useGetSitesQuery etc.
  dedupeOperationSuffix: true,
};

const config: CodegenConfig = {
  generates: {
    // ---- API gateway ----
    'src/client/graphql/types.ts': {
      schema,
      plugins: ['typescript'],
    },
    'src/client/': {
      schema,
      documents: ['src/client/graphql/!(metrics).graphql'],
      preset: 'near-operation-file',
      presetConfig: {
        extension: '.generated.ts',
        baseTypesPath: 'graphql/types.ts',
      },
      plugins: ['typescript-operations', 'typescript-react-apollo'],
      config: hooksConfig,
    },
  },
};

export default config;
