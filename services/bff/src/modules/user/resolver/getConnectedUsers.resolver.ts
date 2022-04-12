import {
    Resolver,
    Query,
    Arg,
    UseMiddleware,
    PubSub,
    PubSubEngine,
    Ctx,
} from "type-graphql";
import { Service } from "typedi";
import { ConnectedUserDto } from "../types";
import { UserService } from "../service";
import { TIME_FILTER } from "../../../constants";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { getHeaders } from "../../../common";

@Service()
@Resolver()
export class GetConnectedUsersResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => ConnectedUserDto)
    @UseMiddleware(Authentication)
    async getConnectedUsers(
        @Arg("orgId") orgId: string,
        @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER,
        @PubSub() pubsub: PubSubEngine,
        @Ctx() ctx: Context
    ): Promise<ConnectedUserDto> {
        const user = this.userService.getConnectedUsers(orgId, getHeaders(ctx));
        pubsub.publish("getConnectedUsers", user);
        return user;
    }
}
