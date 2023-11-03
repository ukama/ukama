/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  ApolloGateway,
  IntrospectAndCompose,
  RemoteGraphQLDataSource,
} from "@apollo/gateway";
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
  INVITATION_PORT,
  MEMBER_PORT,
  METRICS_PORT,
  NETWORK_PORT,
  NODE_PORT,
  ORG_PORT,
  PACKAGE_PORT,
  PLANNING_SERVICE_PORT,
  PLAYGROUND_URL,
  RATE_PORT,
  SIM_PORT,
  SUBSCRIBER_PORT,
  USER_PORT,
} from "../common/configs";
import { HTTP401Error, HTTP500Error, Messages } from "../common/errors";
import { logger } from "../common/logger";
import { THeaders } from "../common/types";
import { parseHeaders } from "../common/utils";
import UserApi from "../user/datasource/user_api";
import { UserResDto, WhoamiDto } from "./../user/resolver/types";
import { configureExpress } from "./configureExpress";

function delay(time: any) {
  return new Promise(resolve => setTimeout(resolve, time));
}
let headers: THeaders = {
  auth: {
    Authorization: "",
    Cookie: "",
  },
  orgId: "",
  userId: "",
  orgName: "",
};

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
        { name: "subscriber", url: `http://localhost:${SUBSCRIBER_PORT}` },
        { name: "sim", url: `http://localhost:${SIM_PORT}` },
        { name: "package", url: `http://localhost:${PACKAGE_PORT}` },
        { name: "rate", url: `http://localhost:${RATE_PORT}` },
        { name: "invitation", url: `http://localhost:${INVITATION_PORT}` },
        { name: "member", url: `http://localhost:${MEMBER_PORT}` },
        { name: "planning", url: `http://localhost:${PLANNING_SERVICE_PORT}` },
      ],
      introspectionHeaders: {
        introspection: "true",
      },
    }),
    buildService({ url }) {
      return new RemoteGraphQLDataSource({
        url,
        willSendRequest({ request }: any) {
          if (request.http.headers.get("introspection") !== "true") {
            request.http.headers.set(
              "x-session-token",
              headers.auth.Authorization
            );
            request.http.headers.set("cookie", headers.auth.Cookie);
            request.http.headers.set("orgId", headers.orgId);
            request.http.headers.set("userId", headers.userId);
            request.http.headers.set("orgName", headers.orgName);
          }
        },
      });
    },
  });
  return gateway;
};

const startServer = async () => {
  await delay(10000);
  const gateway = await loadServers();
  const server = new ApolloServer({
    gateway,
    plugins: [
      ApolloServerPluginInlineTrace({}),
      ApolloServerPluginDrainHttpServer({ httpServer }),
      {
        async requestDidStart(requestContext: any) {
          headers = parseHeaders(requestContext?.request.http.headers);
        },
      },
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
  await new Promise((resolve: any) =>
    httpServer.listen({ port: GATEWAY_PORT }, resolve)
  );

  app.get("/ping", async (_, res) => {
    const r = await fetch(`localhost:${METRICS_PORT}/ping`);
    if (r.status === 200) res.send("pong");
    else res.send(new HTTP500Error("Metrics service ping failed"));
  });
  // const TEMP_KID = "018688fa-d861-4e7b-b119-ffc5e1637ba8";
  // const TEMP_KID = "a9a3dc45-fe06-43d6-b148-7508c9674627";
  app.get("/get-user", async (req, res) => {
    const kId = req.query["kid"] as string;
    if (kId) {
      const userApi = new UserApi();
      const user: UserResDto = await userApi.auth(kId);
      if (user.uuid) {
        // if (TEMP_KID) {
        const whoamiRes: WhoamiDto = await userApi.whoami(user.uuid);
        res.setHeader("Access-Control-Allow-Origin", "http://localhost:4455");
        res.setHeader("Access-Control-Allow-Credentials", "true");
        res.send(whoamiRes);
        return;
      }
    }
    res.send(new HTTP401Error(Messages.HEADER_ERR_USER));
    return;
  });
  logger.info(`🚀 Server ready at http://localhost:${GATEWAY_PORT}/graphql`);
};

startServer();
