import "reflect-metadata";

import { ApolloServer } from "@apollo/server";
import { startStandaloneServer } from "@apollo/server/standalone";
import { buildSubgraphSchema, printSubgraphSchema } from "@apollo/subgraph";
import { GraphQLScalarType } from "graphql";
import { DateTimeResolver } from "graphql-scalars";
import gql from "graphql-tag";
import * as tq from "type-graphql";
import { Context, context } from "./common/context";
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

  const server = new ApolloServer<Context>({
    schema: federatedSchema,
    csrfPrevention: false,
  });

  const { url } = await startStandaloneServer(server, {
    context: async () => context,
    listen: { port: 4041 },
  });

  console.log(`🚀 Ukama Planning Tool Subgraph ready at: ${url}`);
};

app();
