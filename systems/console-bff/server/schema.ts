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
 * CONSOLIDATION-DESIGN §2). Every module's resolvers, plus the planning-tool
 * (Prisma) resolvers, are composed here.
 */
import { GraphQLScalarType, GraphQLSchema } from "graphql";
import { DateTimeResolver } from "graphql-scalars";
import * as tq from "type-graphql";

import billingResolvers from "../billing/resolvers";
import componentResolvers from "../component/resolvers";
import controllerResolvers from "../controller/resolvers";
import healthResolvers from "../health/resolvers";
import initResolvers from "../init/resolver";
import invitationResolvers from "../invitation/resolver";
import memberResolvers from "../member/resolver";
import metricResolvers from "../metric/resolver";
import networkResolvers from "../network/resolvers";
import nodeResolvers from "../node/resolvers";
import notificationResolvers from "../notification/resolvers";
import orgResolvers from "../org/resolver";
import packageResolvers from "../package/resolver";
import paymentResolvers from "../payment/resolver";
import rateResolvers from "../rate/resolver";
import reportResolvers from "../report/resolvers";
import simResolvers from "../sim/resolver";
import siteResolvers from "../site/resolvers";
import softwareResolvers from "../software/resolvers";
import subscriberResolvers from "../subscriber/resolver";
import userResolvers from "../user/resolver";

// Explicitly typed so TS contextually checks each element instead of
// inferring a combined literal type (a 21-array spread otherwise produces a
// "union type too complex to represent" error under ts-node).
const ALL_RESOLVERS: CallableFunction[] = [
  ...orgResolvers,
  ...userResolvers,
  ...networkResolvers,
  ...siteResolvers,
  ...memberResolvers,
  ...invitationResolvers,
  ...nodeResolvers,
  ...packageResolvers,
  ...rateResolvers,
  ...simResolvers,
  ...subscriberResolvers,
  ...controllerResolvers,
  ...healthResolvers,
  ...softwareResolvers,
  ...componentResolvers,
  ...billingResolvers,
  ...paymentResolvers,
  ...reportResolvers,
  ...metricResolvers,
  ...notificationResolvers,
  ...initResolvers,
  // planning-tool is intentionally EXCLUDED for phase 1 (its Prisma client
  // needs a configured PLANNING_TOOL_DB + `prisma generate`). See README
  // "Re-enabling planning-tool" before adding planningResolvers back.
];

export const buildAppSchema = async (): Promise<GraphQLSchema> => {
  return tq.buildSchema({
    resolvers: ALL_RESOLVERS as tq.NonEmptyArray<CallableFunction>,
    scalarsMap: [{ type: GraphQLScalarType, scalar: DateTimeResolver }],
    validate: { forbidUnknownValues: false },
  });
};
