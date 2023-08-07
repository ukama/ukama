import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { parseGatewayHeaders, parseHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { NODE_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import NodeAPI from "./dataSource/node-api";
import resolvers from "./resolvers";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req, res }) => {
      const { cache } = server;
      return {
        headers: parseGatewayHeaders(req.headers),
        dataSources: {
          dataSource: new NodeAPI(),
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
