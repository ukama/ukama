import "reflect-metadata";

import { ApolloServer } from "@apollo/server";
import { startStandaloneServer } from "@apollo/server/standalone";
import { buildSubgraphSchema, printSubgraphSchema } from "@apollo/subgraph";
import { GraphQLScalarType } from "graphql";
import { DateTimeResolver } from "graphql-scalars";
import gql from "graphql-tag";
import * as tq from "type-graphql";
import { Context, context } from "./common/context";
// import { APOptions } from "./common/enums";
import resolvers from "./modules";

const app = async () => {
  // tq.registerEnumType(APOptions, {
  //   name: "APOptions",
  // });
  const ts = await tq.buildSchema({
    resolvers: resolvers,
    scalarsMap: [{ type: GraphQLScalarType, scalar: DateTimeResolver }],
    validate: { forbidUnknownValues: false },
  });

  const federatedSchema = buildSubgraphSchema({
    typeDefs: gql(printSubgraphSchema(ts)),
    resolvers: tq.createResolversMap(ts) as any,
  });

  const server = new ApolloServer<Context>({ schema: federatedSchema });

  const { url } = await startStandaloneServer(server, {
    context: async () => context,
  });

  console.log(`🚀 Ukama Planning Tool Subgraph ready at: ${url}`);
};

app();
