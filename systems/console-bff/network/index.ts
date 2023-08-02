import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import SubGraphServer from "./../common/apollo";
import { NETWORK_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import NetworkAPI from "./datasource/network_api";
import resolvers from "./resolvers";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async () => {
      const { cache } = server;
      return {
        dataSources: {
          dataSource: new NetworkAPI(),
        },
      };
    },
    listen: { port: NETWORK_PORT },
  });

  logger.info(
    `🚀 Ukama Network service running at http://localhost:${NETWORK_PORT}/graphql`
  );
};

runServer();
