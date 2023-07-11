import "dotenv/config";
import * as fs from "fs";
import { createPubSub, createSchema, createYoga } from "graphql-yoga";
import { createServer } from "node:http";
import * as path from "path";

import { METRICS_PORT } from "./../../common/configs";
import resolvers from "./resolvers";

const typeDefs = fs.readFileSync(path.join(process.cwd(), "schema.graphql"), {
  encoding: "utf-8",
});

const pubSub = createPubSub({});

const yoga = createYoga({
  schema: createSchema({
    typeDefs,
    resolvers,
  }),
  logging: true,
  context: {
    pubSub,
  },
});

const server = createServer(yoga);

server.listen(METRICS_PORT, () => {
  console.info("Server is running on http://localhost:4000/graphql");
});
