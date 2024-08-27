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
  BASE_DOMAIN,
  CONSOLE_APP_URL,
  GATEWAY_PORT,
  PLAYGROUND_URL,
  SUBSCRIPTIONS_PORT,
  SUB_GRAPH_LIST,
} from "../common/configs";
import { HTTP401Error, HTTP500Error, Messages } from "../common/errors";
import { logger } from "../common/logger";
import { openStore } from "../common/storage";
import { THeaders } from "../common/types";
import { parseHeaders } from "../common/utils";
import InitAPI from "../init/datasource/init_api";
import { configureExpress } from "./configureExpress";

const COOKIE_EXPIRY_TIME = 3017874138705;

function delay(time: any) {
  return new Promise(resolve => setTimeout(resolve, time));
}
let headers: THeaders = {
  auth: {
    Authorization: "",
    Cookie: "",
  },
  token: "",
  orgId: "",
  userId: "",
  orgName: "",
};

const app = configureExpress(logger);
const httpServer = createServer(app);

const loadServers = async () => {
  const gateway = new ApolloGateway({
    supergraphSdl: new IntrospectAndCompose({
      subgraphs: SUB_GRAPH_LIST.map(({ name, url }: any) => ({
        name,
        url,
      })),

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
            request.http.headers.set("token", headers.token);
          }
        },
      });
    },
  });
  return gateway;
};

const startServer = async () => {
  await delay(10000);
  const store = openStore();
  const gateway = await loadServers();
  const server = new ApolloServer({
    gateway,
    csrfPrevention: true,
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

  app.get("/ping", (_, res) => {
    fetch(`localhost:${SUBSCRIPTIONS_PORT}/ping`)
      .then(r => {
        if (r.status === 200) res.send("pong");
        else res.send(new HTTP500Error("Subscriptions service ping failed"));
      })
      .catch(err => {
        res.send(new HTTP500Error("Subscriptions service ping failed: " + err));
      });
  });

  app.get("/set-theme", (req, res) => {
    const theme = req.query.theme;
    res.cookie("theme", theme, {
      domain: BASE_DOMAIN,
      secure: true,
      sameSite: "lax",
      maxAge: COOKIE_EXPIRY_TIME - (new Date().getTime() - 2017874138705),
      httpOnly: false,
      path: "/",
    });
    res.send("Theme set successfully");
  });

  app.get("/get-user", async (req, res) => {
    const cookies = req.headers["cookie"];
    const initAPI = new InitAPI();
    if (cookies) {
      const sessionRes = await initAPI.validateSession(store, cookies);
      res.setHeader("Content-Type", "application/json");
      res.setHeader("cache-control", "max-age=3600");
      return res.send(sessionRes);
    } else {
      res.send(new HTTP401Error(Messages.HEADER_ERR_USER));
    }
  });

  logger.info(`🚀 Server ready at http://localhost:${GATEWAY_PORT}/graphql`);
};

startServer();
