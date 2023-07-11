import { ApolloServer } from "@apollo/server";
import { ApolloServerPluginInlineTrace } from "@apollo/server/plugin/inlineTrace";
import { startStandaloneServer } from "@apollo/server/standalone";
import { buildSubgraphSchema, printSubgraphSchema } from "@apollo/subgraph";
import { GraphQLScalarType } from "graphql";
import { DateTimeResolver } from "graphql-scalars";
import gql from "graphql-tag";
import "reflect-metadata";
import * as tq from "type-graphql";

import { PLANNING_SERVICE_PORT } from "../../common/configs";
import { logger } from "../../common/logger";
import { PrismaContext, context } from "../../common/prisma";
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

  const { url } = await startStandaloneServer(server, {
    context: async () => context,
    listen: { port: PLANNING_SERVICE_PORT },
  });
  logger.info(
    `🚀 Ukama Planning Tool running at http://localhost:${PLANNING_SERVICE_PORT}/graphql`
  );
};

app();
