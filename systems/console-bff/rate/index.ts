import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { parseGatewayHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { RATE_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import RateAPI from "./datasource/rate_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req, res }) => {
      const { cache } = server;
      return {
        headers: parseGatewayHeaders(req.headers),
        dataSources: {
          dataSource: new RateAPI(),
        },
      };
    },
    listen: { port: RATE_PORT },
  });

  logger.info(
    `🚀 Ukama Rate service running at http://localhost:${RATE_PORT}/graphql`
  );
};

runServer();
