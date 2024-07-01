/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { ApolloServer } from "@apollo/server";
import { expressMiddleware } from "@apollo/server/express4";
import { ApolloServerPluginDrainHttpServer } from "@apollo/server/plugin/drainHttpServer";
import { ApolloServerPluginInlineTrace } from "@apollo/server/plugin/inlineTrace";
import { buildSubgraphSchema, printSubgraphSchema } from "@apollo/subgraph";
import bodyParser from "body-parser";
import cors from "cors";
import express from "express";
import { GraphQLScalarType } from "graphql";
import { DateTimeResolver } from "graphql-scalars";
import gql from "graphql-tag";
import { useServer } from "graphql-ws/lib/use/ws";
import { createServer } from "http";
import "reflect-metadata";
import * as tq from "type-graphql";
import { WebSocketServer } from "ws";

import {
  AUTH_APP_URL,
  CONSOLE_APP_URL,
  PLAYGROUND_URL,
} from "../common/configs";
import { logger } from "../common/logger";
import { storeInStorage } from "../common/storage";
import { SUBSCRIPTIONS_PORT } from "./../common/configs/index";
import resolvers from "./resolvers";

const app = express();
const httpServer = createServer(app);
storeInStorage("UkamaSubscriptions", "running");

const runServer = async () => {
  const ts = await tq.buildSchema({
    resolvers: resolvers,
    scalarsMap: [{ type: GraphQLScalarType, scalar: DateTimeResolver }],
    validate: { forbidUnknownValues: false },
  });

  const wsServer = new WebSocketServer({
    server: httpServer,
    path: "/graphql",
  });

  const federatedSchema = buildSubgraphSchema({
    typeDefs: gql(printSubgraphSchema(ts)),
    resolvers: tq.createResolversMap(ts) as any,
  });

  const serverCleanup = useServer({ schema: federatedSchema }, wsServer);

  const server = new ApolloServer({
    schema: federatedSchema,
    csrfPrevention: false,
    plugins: [
      ApolloServerPluginInlineTrace({}),
      ApolloServerPluginDrainHttpServer({ httpServer }),

      {
        async serverWillStart() {
          return {
            async drainServer() {
              await serverCleanup.dispose();
            },
          };
        },
      },
    ],
  });

  await server.start();
  app.get("/ping", (_, res) => {
    res.send("pong");
  });
  app.use(
    "/graphql",
    cors({
      origin: [AUTH_APP_URL, PLAYGROUND_URL, CONSOLE_APP_URL],
      credentials: true,
    }),
    bodyParser.json(),
    expressMiddleware(server)
  );
  httpServer.listen(SUBSCRIPTIONS_PORT, () => {
    logger.info(
      `🚀 Ukama Subscriptions service running at http://localhost:${SUBSCRIPTIONS_PORT}/graphql`
    );
    logger.info(
      `🚀 Ukama Subscription service running at ws://localhost:${SUBSCRIPTIONS_PORT}/graphql`
    );
  });
};

runServer();
