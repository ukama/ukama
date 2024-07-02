/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { ApolloServer } from "@apollo/server";
import { ApolloServerPluginInlineTrace } from "@apollo/server/plugin/inlineTrace";
import { startStandaloneServer } from "@apollo/server/standalone";
import { buildSubgraphSchema, printSubgraphSchema } from "@apollo/subgraph";
import { GraphQLScalarType } from "graphql";
import { DateTimeResolver } from "graphql-scalars";
import gql from "graphql-tag";
import "reflect-metadata";
import * as tq from "type-graphql";

import { PLANNING_SERVICE_PORT } from "../common/configs";
import { logger } from "../common/logger";
import { PrismaContext, context } from "../common/prisma";
import resolvers from "./modules";

const app = async () => {
  const ts = await tq.buildSchema({
    resolvers: resolvers,
    scalarsMap: [{ type: GraphQLScalarType, scalar: DateTimeResolver }],
    validate: { forbidUnknownValues: false },
  });

  const federatedSchema = buildSubgraphSchema({
    typeDefs: gql(printSubgraphSchema(ts)),
    resolvers: tq.createResolversMap(ts) as any,
  });

  const server = new ApolloServer<PrismaContext>({
    schema: federatedSchema,
    csrfPrevention: false,
    plugins: [ApolloServerPluginInlineTrace({})],
  });

  await startStandaloneServer(server, {
    context: async () => context,
    listen: { port: PLANNING_SERVICE_PORT },
  });
  logger.info(
    `🚀 Ukama Planning Tool running at http://localhost:${PLANNING_SERVICE_PORT}/graphql`
  );
};

app();

// "build": "yarn prisma-pre-build && tsc",
// "all-dev": "concurrently \"yarn notification-dev\" \"yarn init-dev\" \"yarn org-dev\" \"yarn user-dev\" \"yarn network-dev\" \"yarn node-dev\" \"yarn subscriber-dev\" \"yarn sim-dev\" \"yarn package-dev\" \"yarn rate-dev\" \"yarn gateway-dev\" \"yarn invitation-dev\" \"yarn member-dev\" \"yarn planning-tool-dev\"",
// "planning-tool-dev": "prisma generate --schema=./planning-tool/prisma/schema.prisma && nodemon --watch \"planning-tool/**\" --ext \"ts,json,graphql\" --exec \"ts-node planning-tool/index.ts\"",
// "all-start": "node  dist/org/index.js & node dist/notification/index.js & node dist/user/index.js & node dist/network/index.js & node dist/node/index.js & node dist/subscriber/index.js & node dist/sim/index.js & node dist/package/index.js & node dist/rate/index.js & node dist/invitation/index.js & node dist/member/index.js & node dist/planning-tool/index.js & node dist/gateway/index.js"
