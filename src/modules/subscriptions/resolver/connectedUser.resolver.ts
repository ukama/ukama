import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { ConnectedUserDto } from "../../user/types";

@Service()
@Resolver()
export class ConnectedUsersSubscriptionResolver {
    @Subscription(() => ConnectedUserDto, {
        topics: "getConnectedUsers",
    })
    async getConnectedUsers(
        @Root() user: ConnectedUserDto
    ): Promise<ConnectedUserDto> {
        return user;
    }
}
