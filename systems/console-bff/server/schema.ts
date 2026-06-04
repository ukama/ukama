/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Single merged schema for the consolidated API server — plain type-graphql,
 * no federation (the modules share no entity references; see
 * CONSOLIDATION-DESIGN §2). Module resolver arrays are appended here as each
 * Phase B batch migrates.
 */
import { GraphQLScalarType, GraphQLSchema } from "graphql";
import { DateTimeResolver } from "graphql-scalars";
import * as tq from "type-graphql";

import invitationResolvers from "../invitation/resolver";
import memberResolvers from "../member/resolver";
import networkResolvers from "../network/resolvers";
import nodeResolvers from "../node/resolvers";
import orgResolvers from "../org/resolver";
import siteResolvers from "../site/resolvers";
import userResolvers from "../user/resolver";

const ALL_RESOLVERS = [
  ...orgResolvers,
  ...userResolvers,
  ...networkResolvers,
  ...siteResolvers,
  ...memberResolvers,
  ...invitationResolvers,
  ...nodeResolvers,
  // …appended per Phase B batch
] as tq.NonEmptyArray<CallableFunction>;

export const buildAppSchema = async (): Promise<GraphQLSchema> => {
  return tq.buildSchema({
    resolvers: ALL_RESOLVERS,
    scalarsMap: [{ type: GraphQLScalarType, scalar: DateTimeResolver }],
    validate: { forbidUnknownValues: false },
  });
};
