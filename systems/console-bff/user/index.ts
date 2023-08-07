import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { parseGatewayHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { USER_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import UserAPI from "./datasource/user_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req, res }) => {
      const { cache } = server;
      return {
        headers: parseGatewayHeaders(req.headers),
        dataSources: {
          dataSource: new UserAPI(),
        },
      };
    },
    listen: { port: USER_PORT },
  });

  logger.info(
    `🚀 Ukama User service running at http://localhost:${USER_PORT}/graphql`
  );
};

runServer();
