import { ApolloGateway, IntrospectAndCompose } from "@apollo/gateway";
import { ApolloServer } from "@apollo/server";
import { expressMiddleware } from "@apollo/server/express4";
import { ApolloServerPluginDrainHttpServer } from "@apollo/server/plugin/drainHttpServer";
import { ApolloServerPluginInlineTrace } from "@apollo/server/plugin/inlineTrace";
import { json } from "body-parser";
import cors from "cors";
import { createServer } from "http";

import {
  AUTH_APP_URL,
  CONSOLE_APP_URL,
  GATEWAY_PORT,
  NETWORK_PORT,
  NODE_PORT,
  ORG_PORT,
  PLANNING_SERVICE_PORT,
  PLAYGROUND_URL,
  USER_PORT,
} from "../common/configs";
import { logger } from "../common/logger";
import { configureExpress } from "./configureExpress";

function delay(time: any) {
  return new Promise(resolve => setTimeout(resolve, time));
}

const app = configureExpress(logger);
const httpServer = createServer(app);

const loadServers = async () => {
  const gateway = new ApolloGateway({
    supergraphSdl: new IntrospectAndCompose({
      subgraphs: [
        { name: "org", url: `http://localhost:${ORG_PORT}` },
        { name: "node", url: `http://localhost:${NODE_PORT}` },
        { name: "user", url: `http://localhost:${USER_PORT}` },
        { name: "network", url: `http://localhost:${NETWORK_PORT}` },
        { name: "planning", url: `http://localhost:${PLANNING_SERVICE_PORT}` },
      ],
    }),
  });
  return gateway;
};

const startServer = async () => {
  await delay(5000);
  const gateway = await loadServers();
  const server = new ApolloServer({
    gateway,
    plugins: [
      ApolloServerPluginInlineTrace({}),
      ApolloServerPluginDrainHttpServer({ httpServer }),
    ],
  });

  await server.start();

  app.use(
    "/graphql",
    cors({
      origin: [AUTH_APP_URL, PLAYGROUND_URL, CONSOLE_APP_URL],
      credentials: true,
    }),
    json(),
    expressMiddleware(server)
  );
  app.get("/ping", (req, res) => {
    res.send("pong");
  });
  await new Promise((resolve: any) =>
    httpServer.listen({ port: GATEWAY_PORT }, resolve)
  );
  logger.info(`🚀 Server ready at http://localhost:${GATEWAY_PORT}/graphql`);
};

startServer();
