import dotenv from "dotenv";
dotenv.config({ path: ".env" });
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
} from "@apollographql/graphql-playground-html";

const logger = setupLogger("app");

const initializeApp = async () => {
    const app = configureExpress({
        logger,
    });

    const corsOptions = {
        origin: [
            process.env.CONSOLE_APP_URL ?? "",
            process.env.AUTH_APP_URL ?? "",
        ],
        credentials: true,
    };
    logger.info(`CORS ALLOW: ${JSON.stringify(corsOptions)}`);
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
                subscriptionEndpoint: "/graphql",
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

    app.get("/getCookie", (req, res) => {
        if (hasSession(req.headers.cookie || "")) {
            const { expiry = "", id = "" } = req.query;
            if (expiry && id) {
                const date = new Date(expiry as string);
                res.cookie("id", id, {
                    path: "/",
                    secure: true,
                    expires: date,
                    httpOnly: true,
                    sameSite: "lax",
                    domain: process.env.DOMAIN,
                });
                res.send({ success: true });
            }
        }
        res.send({ success: false });
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
