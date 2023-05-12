import {
    MiddlewareOptions,
    renderPlaygroundPage,
} from "@apollographql/graphql-playground-html";
import cookieParser from "cookie-parser";
import cors from "cors";
import dotenv from "dotenv";
dotenv.config({ path: ".env" });
import { createServer } from "http";
import "reflect-metadata";
import configureApolloServer from "./config/configureApolloServer";
import configureExpress from "./config/configureExpress";
import setupLogger from "./config/setupLogger";
import { PORT } from "./constants";
import { job } from "./jobs/subscriptionJob";
import { mockServer } from "./mockServer";

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

    app.use(cors(corsOptions));
    app.use(cookieParser());

    const { server, schema } = await configureApolloServer();
    server.applyMiddleware({ app, cors: corsOptions });

    app.use((req, res, next) => {
        res.header("Access-Control-Allow-Origin", [
            process.env.CONSOLE_APP_URL ?? "",
            process.env.AUTH_APP_URL ?? "",
        ]);
        res.header(
            "Access-Control-Allow-Headers",
            "Origin, X-Requested-With, Content-Type, Accept"
        );
        res.setHeader("Access-Control-Allow-Credentials", "true");
        res.header(
            "Access-Control-Allow-Methods",
            "GET, POST, PUT, DELETE, OPTIONS"
        );
        next();
    });

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
                res.setHeader(
                    "Access-Control-Allow-Origin",
                    process.env.AUTH_APP_URL ?? ""
                );
                res.setHeader("Access-Control-Allow-Credentials", "true");
                res.send({ success: true });
                res.end();
            }
        } else {
            res.send({ success: false });
            res.end();
        }
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
