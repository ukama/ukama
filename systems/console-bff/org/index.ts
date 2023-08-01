import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import SubGraphServer from "./../common/apollo";
import { ORG_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import OrgAPI from "./datasource/org_api";
import UserAPI from "./datasource/user_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async () => {
      const { cache } = server;
      return {
        dataSources: {
          dataSource: new OrgAPI(),
          dataSoureceUser: new UserAPI(),
        },
      };
    },
    listen: { port: ORG_PORT },
  });

  logger.info(
    `🚀 Ukama Node service running at http://localhost:${ORG_PORT}/graphql`
  );
};

runServer();
