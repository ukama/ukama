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
        if (!context.req.headers.authorization)
            throw new HTTP401Error(Messages.ERR_REQUIRED_HEADER_NOT_FOUND);

        const session = context.req.headers.authorization;

        context.session = session;
        return next();
    }
}
