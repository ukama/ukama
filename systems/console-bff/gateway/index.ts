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
import { RootDatabase } from "lmdb";

import {
  AUTH_APP_URL,
  BASE_DOMAIN,
  CONSOLE_APP_URL,
  GATEWAY_PORT,
  INTROSPECTION_ENABLED,
  PLAYGROUND_URL,
  SUBSCRIPTIONS_PORT,
  SUB_GRAPH_LIST,
} from "../common/configs";
import { validateEnv } from "../common/configs/validateEnv";
import { HTTP401Error, HTTP500Error, Messages } from "../common/errors";
import { logger } from "../common/logger";
import { closeStore, openStore } from "../common/storage";
import { THeaders } from "../common/types";
import { parseExpressHeaders } from "../common/utils";
import InitAPI from "../init/datasource/init_api";
import { configureExpress } from "./configureExpress";

const COOKIE_EXPIRY_TIME = 3017874138705;

// Request body cap — GraphQL operations are small; this blocks payload DoS.
const JSON_BODY_LIMIT = "1mb";

// Gateway composition retry/backoff so a subgraph that is briefly down at
// boot doesn't permanently fail startup.
const COMPOSE_MAX_RETRIES = 30;
const COMPOSE_BASE_DELAY_MS = 2000;
const COMPOSE_MAX_DELAY_MS = 30_000;

interface GatewayContext {
  headers: THeaders;
  requestId: string;
}

/** Readiness flag flipped true once the gateway has composed and is listening. */
let isReady = false;

const EMPTY_HEADERS: THeaders = {
  auth: { Authorization: "", Cookie: "" },
  token: "",
  orgId: "",
  userId: "",
  orgName: "",
};

/**
 * True if the parsed request body is a GraphQL introspection query. Used to
 * let schema-only tooling (codegen) through without a session. Matches the
 * standard introspection operation name or the `__schema` meta-field.
 */
const isIntrospectionRequest = (body: unknown): boolean => {
  if (!body || typeof body !== "object") return false;
  const ops = Array.isArray(body) ? body : [body];
  return ops.some(op => {
    const o = op as { operationName?: string; query?: string };
    return (
      o?.operationName === "IntrospectionQuery" ||
      (typeof o?.query === "string" && o.query.includes("__schema"))
    );
  });
};

function delay(time: number) {
  return new Promise(resolve => setTimeout(resolve, time));
}

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
      return new RemoteGraphQLDataSource<GatewayContext>({
        url,
        willSendRequest({ request, context }) {
          // Auth headers are read from the per-request context so
          // concurrent requests can never leak each other's identity.
          if (request.http?.headers.get("introspection") === "true") return;
          // Propagate the correlation id to the subgraph for tracing.
          if (context?.requestId) {
            request.http?.headers.set("x-request-id", context.requestId);
          }
          const requestHeaders = context?.headers;
          if (!requestHeaders) return;
          request.http?.headers.set(
            "x-session-token",
            requestHeaders.auth.Authorization
          );
          request.http?.headers.set("cookie", requestHeaders.auth.Cookie);
          request.http?.headers.set("token", requestHeaders.token);
        },
      });
    },
  });
  return gateway;
};

/**
 * Installs SIGTERM/SIGINT handlers that stop accepting new work, drain
 * in-flight requests, and close the Apollo server, HTTP server and LMDB
 * store before exiting — so rolling restarts don't drop requests.
 */
const registerShutdown = (
  server: ApolloServer<GatewayContext>,
  store: RootDatabase
): void => {
  let shuttingDown = false;
  const shutdown = async (signal: string) => {
    if (shuttingDown) return;
    shuttingDown = true;
    isReady = false; // fail readiness so traffic drains away
    logger.info(`Received ${signal}, shutting down gracefully...`);
    try {
      await server.stop(); // drains via ApolloServerPluginDrainHttpServer
      await new Promise<void>(resolve => httpServer.close(() => resolve()));
      await closeStore(store);
      logger.info("Shutdown complete");
      process.exit(0);
    } catch (err) {
      logger.error(`Error during shutdown: ${err}`);
      process.exit(1);
    }
  };
  process.on("SIGTERM", () => void shutdown("SIGTERM"));
  process.on("SIGINT", () => void shutdown("SIGINT"));
};

/**
 * Composes the federated gateway and starts the Apollo server, retrying with
 * exponential backoff so a subgraph that is briefly unavailable at boot does
 * not permanently fail startup.
 */
const startApolloWithRetry = async (): Promise<
  ApolloServer<GatewayContext>
> => {
  let attempt = 0;
  // eslint-disable-next-line no-constant-condition
  while (true) {
    attempt += 1;
    try {
      const gateway = await loadServers();
      const server = new ApolloServer<GatewayContext>({
        gateway,
        csrfPrevention: true,
        // Off in production unless ENABLE_INTROSPECTION=true (e.g. for codegen).
        introspection: INTROSPECTION_ENABLED,
        plugins: [
          ApolloServerPluginInlineTrace({}),
          ApolloServerPluginDrainHttpServer({ httpServer }),
        ],
      });
      await server.start();
      logger.info(`Gateway composed on attempt ${attempt}`);
      return server;
    } catch (err) {
      if (attempt >= COMPOSE_MAX_RETRIES) {
        throw new Error(
          `Gateway composition failed after ${attempt} attempts: ${err}`
        );
      }
      const backoff = Math.min(
        COMPOSE_BASE_DELAY_MS * 2 ** (attempt - 1),
        COMPOSE_MAX_DELAY_MS
      );
      logger.warn(
        `Gateway composition attempt ${attempt} failed (a subgraph may be down); retrying in ${backoff}ms: ${err}`
      );
      await delay(backoff);
    }
  }
};

const startServer = async () => {
  validateEnv();
  const store = openStore();

  // Liveness: the process is up. Registered before composition so the
  // container is reported alive while subgraphs are still coming online.
  app.get("/healthz", (_, res) => {
    res.status(200).json({ status: "ok" });
  });

  // Readiness: only 200 once the gateway has composed and is listening.
  app.get("/readyz", (_, res) => {
    if (isReady) res.status(200).json({ status: "ready" });
    else res.status(503).json({ status: "not-ready" });
  });

  await delay(10000);
  const server = await startApolloWithRetry();

  app.use(
    "/graphql",
    cors({
      origin: [AUTH_APP_URL, PLAYGROUND_URL, CONSOLE_APP_URL],
      credentials: true,
    }),
    json({ limit: JSON_BODY_LIMIT }),
    expressMiddleware(server, {
      context: async ({ req }) => {
        const requestId = (req.headers["x-request-id"] as string) ?? "";
        // Schema introspection (e.g. codegen) carries no session. Allow it
        // past the auth gate when introspection is enabled — it reveals only
        // the schema, never data — so tooling can read the schema.
        if (INTROSPECTION_ENABLED && isIntrospectionRequest(req.body)) {
          return { headers: EMPTY_HEADERS, requestId };
        }
        return {
          headers: parseExpressHeaders(req.headers),
          requestId,
        };
      },
    })
  );

  await new Promise((resolve: any) =>
    httpServer.listen({ port: GATEWAY_PORT }, resolve)
  );
  isReady = true;

  registerShutdown(server, store);

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
      secure: false,
      sameSite: "lax",
      maxAge: COOKIE_EXPIRY_TIME - (new Date().getTime() - 2017874138705),
      httpOnly: false,
      path: "/",
    });
    res.send("Theme set successfully");
  });

  app.get("/get-user", async (req, res) => {
    const cookies = req.headers["cookie"];
    if (!cookies) {
      return res.status(401).send(new HTTP401Error(Messages.HEADER_ERR_USER));
    }
    try {
      const initAPI = new InitAPI();
      const sessionRes = await initAPI.validateSession(store, cookies);
      res.setHeader("Content-Type", "application/json");
      return res.send(sessionRes);
    } catch (err) {
      logger.error(`get-user failed: ${err}`);
      return res
        .status(500)
        .send(new HTTP500Error("Failed to validate session"));
    }
  });

  logger.info(`🚀 Server ready at http://localhost:${GATEWAY_PORT}/graphql`);
};

process.on("unhandledRejection", reason => {
  logger.error(`Unhandled promise rejection: ${reason}`);
});

process.on("uncaughtException", err => {
  logger.error(`Uncaught exception: ${err.stack ?? err}`);
  process.exit(1);
});

startServer().catch(err => {
  logger.error(`Gateway failed to start: ${err.stack ?? err}`);
  process.exit(1);
});
