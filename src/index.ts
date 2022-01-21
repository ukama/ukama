import dotenv from "dotenv";
dotenv.config({ path: `.env.${process.env.NODE_ENV}` });

import cors from "cors";
import "reflect-metadata";
import { PORT } from "./constants";
import { createServer } from "http";
import cookieParser from "cookie-parser";
import { mockServer } from "./mockServer";
import { job } from "./jobs/subscriptionJob";
import setupLogger from "./config/setupLogger";
import configureExpress from "./config/configureExpress";
import configureApolloServer from "./config/configureApolloServer";
import {
    MiddlewareOptions,
    renderPlaygroundPage,
    RenderPageOptions as PlaygroundRenderPageOptions,
} from "@apollographql/graphql-playground-html";

const logger = setupLogger("app");

const initializeApp = async () => {
    const app = configureExpress({
        logger,
    });

    const corsOptions = {
        origin: process.env.DASHBOARD_APP_URL,
        credentials: true,
    };
    app.use(cors(corsOptions));
    app.use(cookieParser());

    const { server, schema } = await configureApolloServer();
    server.applyMiddleware({ app, cors: false });

    const httpServer = createServer(app);
    server.installSubscriptionHandlers(httpServer);

    app.get("/playground", (req, res) => {
        if (hasSession(req.headers.cookie || "")) {
            const playgroundRenderPageOptions: MiddlewareOptions = {
                endpoint: "/graphql",
            };

            const playground = renderPlaygroundPage(
                playgroundRenderPageOptions
            );
            res.write(playground);
            res.end();
        } else {
            res.redirect(
                `${process.env.AUTH_APP_URL}?redirect=${req.headers.host}${req.originalUrl}`
            );
            res.end();
        }
    });

    app.get("/ping", (req, res) => {
        res.send("pong");
    });

    mockServer(app);
    httpServer.listen(PORT, () => {
        logger.info(`Server listening on port: ${PORT}`);
        job(schema);
    });
};

initializeApp().catch(error => logger.error(error));

const hasSession = (session: string) => {
    if (session) {
        return session.includes("ukama_session");
    }
    return false;
};
