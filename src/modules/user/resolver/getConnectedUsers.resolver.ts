import {
    Resolver,
    Query,
    Arg,
    UseMiddleware,
    PubSub,
    PubSubEngine,
} from "type-graphql";
import { Service } from "typedi";
import { ConnectedUserDto } from "../types";
import { UserService } from "../service";
import { TIME_FILTER } from "../../../constants";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetConnectedUsersResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => ConnectedUserDto)
    @UseMiddleware(Authentication)
    async getConnectedUsers(
        @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<ConnectedUserDto> {
        const user = this.userService.getConnectedUsers(filter);
        pubsub.publish("getConnectedUsers", user);
        return user;
    }
}
