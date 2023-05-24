import {
    Arg,
    Ctx,
    PubSub,
    PubSubEngine,
    Query,
    Resolver,
    UseMiddleware,
} from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { TIME_FILTER } from "../../../constants";
import { UserService } from "../service";
import { ConnectedUserDto } from "../types";

@Service()
@Resolver()
export class GetConnectedUsersResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => ConnectedUserDto)
    @UseMiddleware(Authentication)
    async getConnectedUsers(
        @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER,
        @PubSub() pubsub: PubSubEngine,
        @Ctx() ctx: Context
    ): Promise<ConnectedUserDto> {
        const user = this.userService.getConnectedUsers(parseHeaders(ctx));
        pubsub.publish("getConnectedUsers", user);
        return user;
    }
}
