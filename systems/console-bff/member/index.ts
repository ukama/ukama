import { startStandaloneServer } from "@apollo/server/standalone";
import "reflect-metadata";

import { parseGatewayHeaders } from "../common/utils";
import SubGraphServer from "./../common/apollo";
import { MEMBER_PORT } from "./../common/configs";
import { logger } from "./../common/logger";
import MemberAPI from "./datasource/member_api";
import resolvers from "./resolver";

const runServer = async () => {
  const server = await SubGraphServer(resolvers);
  await startStandaloneServer(server, {
    context: async ({ req }) => {
      return {
        headers: parseGatewayHeaders(req.headers),
        dataSources: {
          dataSource: new MemberAPI(),
        },
      };
    },
    listen: { port: MEMBER_PORT },
  });

  logger.info(
    `🚀 Ukama Member service running at http://localhost:${MEMBER_PORT}/graphql`
  );
};

runServer();
