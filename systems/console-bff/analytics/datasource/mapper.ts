/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * The analytics gateway renders protobuf responses with protojson
 * (EmitUnpopulated), so the wire JSON already uses lowerCamelCase keys,
 * RFC3339 string timestamps and emits zero values for every field. That means
 * the parsed response is shape-compatible with the GraphQL ObjectTypes —
 * unlike the snake_case REST DTOs elsewhere, no per-field remapping is needed.
 *
 * This single typed passthrough is the one place to add reshaping later if the
 * gateway contract ever diverges from the GraphQL types.
 */
export const mapAnalytics = <T>(res: unknown): T => res as T;
