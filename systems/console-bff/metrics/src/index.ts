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

import { METRICS_PORT } from "../../common/configs";
import { logger } from "../../common/logger";
import resolvers from "./resolvers";

const app = express();
const httpServer = createServer(app);

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
  app.use(
    "/graphql",
    cors<cors.CorsRequest>(),
    bodyParser.json(),
    expressMiddleware(server)
  );
  httpServer.listen(METRICS_PORT, () => {
    logger.info(
      `🚀 Ukama Metrics service running at http://localhost:${METRICS_PORT}/graphql`
    );
    logger.info(
      `🚀 Ukama Metrics subscription service running at ws://localhost:${METRICS_PORT}/graphql`
    );
  });
};

runServer();
