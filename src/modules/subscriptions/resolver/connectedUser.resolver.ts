import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { ConnectedUserDto } from "../../user/types";

@Service()
@Resolver()
export class ConnectedUsersSubscriptionResolver {
    @Subscription(() => ConnectedUserDto, {
        topics: "CONNECTEDUSER",
    })
    async CONNECTEDUSER(
        @Root() user: ConnectedUserDto
    ): Promise<ConnectedUserDto> {
        return user;
    }
}
