import { ApolloServer } from "@apollo/server";
import { ApolloServerPluginInlineTrace } from "@apollo/server/plugin/inlineTrace";
import { startStandaloneServer } from "@apollo/server/standalone";
import { buildSubgraphSchema, printSubgraphSchema } from "@apollo/subgraph";
import express from "express";
import { GraphQLScalarType } from "graphql";
import { DateTimeResolver } from "graphql-scalars";
import gql from "graphql-tag";
import "reflect-metadata";
import * as tq from "type-graphql";

import { NODE_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import NodeAPI from "./dataSource/node-api";
import resolvers from "./resolvers";

const app = express();
const runServer = async () => {
  const ts = await tq.buildSchema({
    resolvers: resolvers,
    scalarsMap: [{ type: GraphQLScalarType, scalar: DateTimeResolver }],
    validate: { forbidUnknownValues: false },
  });

  const federatedSchema = buildSubgraphSchema({
    typeDefs: gql(printSubgraphSchema(ts)),
    resolvers: tq.createResolversMap(ts) as any,
  });

  const server = new ApolloServer({
    schema: federatedSchema,
    csrfPrevention: false,

    plugins: [ApolloServerPluginInlineTrace({})],
  });

  await startStandaloneServer(server, {
    context: async () => {
      const { cache } = server;
      return {
        // We create new instances of our data sources with each request,
        // passing in our server's cache.
        dataSources: {
          nodeAPI: new NodeAPI({ cache }),
        },
      };
    },
    listen: { port: NODE_PORT },
  });

  logger.info(
    `🚀 Ukama Node service running at http://localhost:${NODE_PORT}/graphql`
  );
};

runServer();
