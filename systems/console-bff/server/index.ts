/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Consolidated console-bff API server (CONSOLIDATION-DESIGN). Replaces the
 * federation gateway + 22 subgraph processes with a single Apollo server over
 * one merged schema. The subscriptions service stays separate. Phase A: runs
 * on API_PORT (8090) alongside the legacy gateway; becomes :8080 at cutover.
 */
import { ApolloServer } from "@apollo/server";
import { ApolloServerPluginDrainHttpServer } from "@apollo/server/plugin/drainHttpServer";
import {
  ApolloServerPluginLandingPageLocalDefault,
  ApolloServerPluginLandingPageProductionDefault,
} from "@apollo/server/plugin/landingPage/default";
// Apollo Server v5: the Express integration lives in its own package.
import { expressMiddleware } from "@as-integrations/express4";
import cors from "cors";
import express from "express";
import { createServer } from "http";
import { RootDatabase } from "lmdb";

import {
  API_PORT,
  AUTH_APP_URL,
  BASE_DOMAIN,
  CONSOLE_APP_URL,
  INTROSPECTION_ENABLED,
  PLAYGROUND_URL,
  SUBSCRIPTIONS_PORT,
} from "../common/configs";
import { validateEnv } from "../common/configs/validateEnv";
import { HTTP401Error, HTTP500Error, Messages } from "../common/errors";
import { logger } from "../common/logger";
import { configureExpress } from "../common/middleware/expressApp";
import { persistedOperations } from "../common/middleware/persistedOperations";
import { addInStore, closeStore, openStore } from "../common/storage";
import { THeaders } from "../common/types";
import { parseExpressHeaders, parseToken } from "../common/utils";
import { ServiceUrlResolver } from "../dashboard/baseUrls";
import InitAPI, {
  SessionValidationError,
  getWelcomeStoreKey,
} from "../init/datasource/init_api";
import { AppContext, buildDataSources, buildHeaders } from "./context";
import { buildAppSchema } from "./schema";

const COOKIE_EXPIRY_TIME = 3017874138705;
const JSON_BODY_LIMIT = "1mb";

let isReady = false;

const EMPTY_HEADERS: THeaders = {
  auth: { Authorization: "", Cookie: "" },
  token: "",
  orgId: "",
  userId: "",
  orgName: "",
};

/** True if the parsed body is a GraphQL introspection query (codegen). */
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

const app = configureExpress(logger);
const httpServer = createServer(app);

const registerShutdown = (
  server: ApolloServer<AppContext>,
  store: RootDatabase
): void => {
  let shuttingDown = false;
  const shutdown = async (signal: string) => {
    if (shuttingDown) return;
    shuttingDown = true;
    isReady = false;
    logger.info(`Received ${signal}, shutting down gracefully...`);
    try {
      await server.stop();
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

const startServer = async () => {
  validateEnv();
  const store = openStore();

  app.get("/healthz", (_, res) => {
    res.status(200).json({ status: "ok" });
  });
  app.get("/readyz", (_, res) => {
    if (isReady) res.status(200).json({ status: "ready" });
    else res.status(503).json({ status: "not-ready" });
  });

  const schema = await buildAppSchema();
  // NODE_ENV=production (Docker) selects Apollo's minimal landing page by
  // default even when introspection is on. Use Sandbox when introspection is
  // explicitly enabled (local docker-compose / codegen).
  const landingPagePlugin = INTROSPECTION_ENABLED
    ? ApolloServerPluginLandingPageLocalDefault({ embed: true })
    : ApolloServerPluginLandingPageProductionDefault();
  const server = new ApolloServer<AppContext>({
    schema,
    csrfPrevention: true,
    introspection: INTROSPECTION_ENABLED,
    plugins: [
      landingPagePlugin,
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
    express.json({ limit: JSON_BODY_LIMIT }),
    // Production allowlist: only operations shipped by the console are
    // accepted when PERSISTED_OPS_ENFORCED=true (no-op otherwise).
    persistedOperations({ allowIntrospection: INTROSPECTION_ENABLED }),
    expressMiddleware(server, {
      context: async ({ req }): Promise<AppContext> => {
        const requestId = (req.headers["x-request-id"] as string) ?? "";
        const dataSources = buildDataSources();
        // Introspection (codegen) carries no session: allow it past the auth
        // gate when introspection is enabled (schema only, never data).
        if (INTROSPECTION_ENABLED && isIntrospectionRequest(req.body)) {
          return {
            headers: EMPTY_HEADERS,
            requestId,
            dataSources,
            urls: new ServiceUrlResolver(""),
          };
        }
        const headers = buildHeaders(req.headers);
        return {
          headers,
          requestId,
          dataSources,
          urls: new ServiceUrlResolver(headers.orgName),
        };
      },
    })
  );

  await new Promise((resolve: any) =>
    httpServer.listen({ port: API_PORT }, resolve)
  );
  isReady = true;
  registerShutdown(server, store);

  // Liveness ping kept as an alias; also checks the subscriptions service.
  app.get("/ping", (_, res) => {
    fetch(`http://localhost:${SUBSCRIPTIONS_PORT}/ping`)
      .then(r => res.send(r.status === 200 ? "pong" : "subscriptions down"))
      .catch(err =>
        res.send(new HTTP500Error("Subscriptions service ping failed: " + err))
      );
  });

  app.get("/set-theme", (req, res) => {
    res.cookie("theme", req.query.theme, {
      domain: BASE_DOMAIN,
      secure: false,
      sameSite: "lax",
      maxAge: COOKIE_EXPIRY_TIME - (new Date().getTime() - 2017874138705),
      httpOnly: false,
      path: "/",
    });
    res.send("Theme set successfully");
  });

  // Marks the welcome page as acknowledged for the authenticated user.
  // Auth: same signed-token verification as /graphql (cookie or
  // x-session-token). Once recorded, /get-user mints subsequent tokens with
  // isShowWelcome=false, so the console stops routing the user to /welcome.
  app.post(
    "/welcome-seen",
    cors({
      origin: [AUTH_APP_URL, PLAYGROUND_URL, CONSOLE_APP_URL],
      credentials: true,
    }),
    async (req, res) => {
      try {
        const headers = parseExpressHeaders(req.headers);
        const userId = parseToken(headers.token, "userId");
        if (!userId) {
          return res
            .status(401)
            .send(new HTTP401Error(Messages.HEADER_ERR_USER));
        }
        await addInStore(store, getWelcomeStoreKey(userId), true);
        return res.status(200).json({ ok: true });
      } catch (err) {
        const reason = err instanceof Error ? err.message : String(err);
        logger.error(`welcome-seen failed: ${reason}`);
        return res.status(401).send(new HTTP401Error(Messages.HEADER_ERR_AUTH));
      }
    }
  );

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
      const reason = err instanceof Error ? err.message : String(err);
      const step = err instanceof SessionValidationError ? err.step : "UNKNOWN";
      logger.error(`get-user failed at ${step}: ${reason}`);
      // 401 (not 500): the session/user mapping is invalid — clients should
      // re-authenticate or land on /unauthorized. Step + reason included so
      // both the console and the logs say which hop broke.
      return res.status(401).json({ step, error: reason });
    }
  });

  logger.info(
    `🚀 Console-BFF API ready at http://localhost:${API_PORT}/graphql (introspection: ${INTROSPECTION_ENABLED})`
  );
};

process.on("unhandledRejection", reason => {
  logger.error(`Unhandled promise rejection: ${reason}`);
});
process.on("uncaughtException", err => {
  logger.error(`Uncaught exception: ${err.stack ?? err}`);
  process.exit(1);
});

startServer().catch(err => {
  logger.error(`API server failed to start: ${err.stack ?? err}`);
  process.exit(1);
});
