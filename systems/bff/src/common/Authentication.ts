import {
    MiddlewareInterface,
    NextFn,
} from "type-graphql/dist/interfaces/Middleware";
import { ResolverData } from "type-graphql/dist/interfaces/ResolverData";
import { Service } from "typedi";
import { HTTP401Error, Messages } from "../errors";
import { Context } from "./types";

@Service()
export class Authentication implements MiddlewareInterface<Context> {
    async use({ context }: ResolverData<Context>, next: NextFn): Promise<void> {
        if (context.req.headers !== undefined) {
            const token = context.req.headers["x-session-token"] || "";
            const cookie =
                context.req.headers.cookie &&
                context.req.headers.cookie.includes("ukama_session")
                    ? context.req.headers.cookie
                    : "";

            if (!cookie && !token) {
                throw new HTTP401Error(Messages.REQUEST_AUTHENTICATION_FAILED);
            }
            context.authType = token ? "token" : "cookie";
        }
        return next();
    }
}
