/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Field, ObjectType, registerEnumType } from "type-graphql";

/**
 * Errors-as-data contract for composite (view-domain) queries — see
 * docs/bff-screen-api-plan.md §4.5. Frontend branches on `code`,
 * never on `message` text.
 */
export enum SectionErrorCode {
  UPSTREAM_TIMEOUT = "UPSTREAM_TIMEOUT",
  UPSTREAM_ERROR = "UPSTREAM_ERROR",
  NOT_FOUND = "NOT_FOUND",
  FORBIDDEN = "FORBIDDEN",
  INTERNAL = "INTERNAL",
  /** Schema is design-complete but the backend endpoint/property does not
   *  exist yet (TODO(backend-gap) in the resolver). The section resolves to
   *  null until the upstream ships; the console renders a placeholder and
   *  needs no change when data arrives. */
  NOT_IMPLEMENTED = "NOT_IMPLEMENTED",
}

registerEnumType(SectionErrorCode, {
  name: "SectionErrorCode",
  description:
    "Machine-readable failure code for a composite query section. " +
    "UI branches on this code; `message` is for display/logs only.",
});

@ObjectType({
  description:
    "Typed failure of one section of a composite query. The section's data " +
    "field resolves to null and a SectionError describes why, so the UI can " +
    "distinguish 'failed' from 'genuinely empty'.",
})
export class SectionError {
  @Field(() => String)
  section: string;

  @Field(() => SectionErrorCode)
  code: SectionErrorCode;

  @Field(() => String)
  message: string;
}

@ObjectType({
  description:
    "Expected/business failure of a mutation (validation, state conflict). " +
    "Rendered inline on forms; unexpected failures still throw.",
})
export class UserError {
  @Field(() => String, { nullable: true })
  field?: string;

  @Field(() => String)
  code: string;

  @Field(() => String)
  message: string;
}
