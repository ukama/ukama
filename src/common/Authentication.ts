import {
    MiddlewareInterface,
    NextFn,
} from "type-graphql/dist/interfaces/Middleware";
import { ResolverData } from "type-graphql/dist/interfaces/ResolverData";
import { Context } from "./types";
import { Service } from "typedi";

@Service()
export class Authentication implements MiddlewareInterface<Context> {
    async use({ context }: ResolverData<Context>, next: NextFn): Promise<void> {
        if (context.req.headers.cookie)
            context.cookie = context.req.headers.cookie;
        return next();
    }
}
