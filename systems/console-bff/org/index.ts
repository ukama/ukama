import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { parseGatewayHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { ORG_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import OrgAPI from "./datasource/org_api";
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
          dataSource: new OrgAPI(),
          dataSoureceUser: new UserAPI(),
        },
      };
    },
    listen: { port: ORG_PORT },
  });

  logger.info(
    `🚀 Ukama Org service running at http://localhost:${ORG_PORT}/graphql`
  );
};

runServer();
