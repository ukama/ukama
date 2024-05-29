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

import { whoami } from "../common/auth/authCalls";
import {
  AUTH_APP_URL,
  BASE_DOMAIN,
  CONSOLE_APP_URL,
  GATEWAY_PORT,
  METRICS_PORT,
  PLAYGROUND_URL,
  SUB_GRAPH_LIST,
} from "../common/configs";
import { HTTP401Error, HTTP500Error, Messages } from "../common/errors";
import { logger } from "../common/logger";
import { THeaders } from "../common/types";
import { parseHeaders } from "../common/utils";
import UserApi from "../user/datasource/user_api";
import { UserResDto, WhoamiDto } from "../user/resolver/types";
import { configureExpress } from "./configureExpress";

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
    fetch(`localhost:${METRICS_PORT}/ping`)
      .then(r => {
        if (r.status === 200) res.send("pong");
        else res.send(new HTTP500Error("Metrics service ping failed"));
      })
      .catch(err => {
        res.send(new HTTP500Error("Metrics service ping failed: " + err));
      });
  });

  app.get("/get-user", async (req, res) => {
    const cookies = req.headers["cookie"];
    let kId = "";
    if (cookies) {
      const whoamiRes = await whoami(cookies);
      if (whoamiRes?.data) {
        kId = whoamiRes.data.identity.id;
      }
    } else {
      res.send(new HTTP401Error(Messages.HEADER_ERR_USER));
    }
    if (kId) {
      const userApi = new UserApi();
      userApi
        .auth(kId)
        .then((user: UserResDto) => {
          if (user.uuid) {
            userApi
              .whoami(user.uuid)
              .then((whoamiRes: WhoamiDto) => {
                const orgId =
                  whoamiRes.ownerOf.length > 0
                    ? whoamiRes.ownerOf[0].id
                    : whoamiRes.memberOf.length > 0
                    ? whoamiRes.memberOf[0].id
                    : "";

                const orgName =
                  whoamiRes.ownerOf.length > 0
                    ? whoamiRes.ownerOf[0].name
                    : whoamiRes.memberOf.length > 0
                    ? whoamiRes.memberOf[0].name
                    : "";

                const cookie = `uid=${user.uuid}&org-id=${orgId}&org-name=${orgName}`;

                const base64Cookie = Buffer.from(cookie).toString("base64");
                res.cookie("user_session", base64Cookie, {
                  domain: BASE_DOMAIN,
                  secure: true,
                  sameSite: "lax",
                  maxAge: 86400,
                  httpOnly: true,
                  path: "/",
                });
                res.setHeader("Access-Control-Allow-Origin", CONSOLE_APP_URL);
                res.setHeader("Access-Control-Allow-Credentials", "true");
                res.send(whoamiRes);
              })
              .catch(err => {
                logger.error(err);
                res.send(new HTTP500Error("Failed to get user details"));
              });
          } else {
            res.send(new HTTP401Error(Messages.HEADER_ERR_USER));
          }
        })
        .catch(err => {
          logger.error(err);
          res.send(new HTTP500Error("Failed to authenticate user"));
        });
    } else {
      res.send(new HTTP401Error(Messages.HEADER_ERR_USER));
    }
  });

  logger.info(`🚀 Server ready at http://localhost:${GATEWAY_PORT}/graphql`);
};

startServer();
