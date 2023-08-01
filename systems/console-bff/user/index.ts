import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import SubGraphServer from "./../common/apollo";
import { USER_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import UserAPI from "./datasource/user_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async () => {
      const { cache } = server;
      return {
        dataSources: {
          dataSource: new UserAPI(),
        },
      };
    },
    listen: { port: USER_PORT },
  });

  logger.info(
    `🚀 Ukama Node service running at http://localhost:${USER_PORT}/graphql`
  );
};

runServer();
