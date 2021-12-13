import {
    MiddlewareInterface,
    NextFn,
} from "type-graphql/dist/interfaces/Middleware";
import { ResolverData } from "type-graphql/dist/interfaces/ResolverData";
import { Context } from "./types";
import { Service } from "typedi";
import { HTTP401Error, Messages } from "../errors";

@Service()
export class Authentication implements MiddlewareInterface<Context> {
    async use({ context }: ResolverData<Context>, next: NextFn): Promise<void> {
        if (!context.req.headers.cookie)
            throw new HTTP401Error(Messages.ERR_REQUIRED_HEADER_NOT_FOUND);

        context.cookie = context.req.headers.cookie;
        return next();
    }
}
